package sim

import (
	"fmt"
	"sort"
)

type AccessType string

const (
	AccessRead  AccessType = "r"
	AccessWrite AccessType = "w"
)

type Perm struct {
	Read  bool
	Write bool
}

type pte struct {
	Frame   int
	Perm    Perm
	Present bool
}

type frameOwner struct {
	PID int
	VPN uint64
}

type tlbEntry struct {
	Valid bool
	PID   int
	VPN   uint64
	Frame int
	Perm  Perm
}

type MemoryManager struct {
	pageSize    uint64
	totalFrames int
	freeFrames  []int
	frames      map[int]frameOwner
	residentQ   []int
	pageTables  map[int]map[uint64]pte
	tlb         []tlbEntry
	tlbNext     int
	faults      FaultCounters
}

func NewMemoryManager(totalFrames, tlbEntries int) *MemoryManager {
	if totalFrames <= 0 {
		totalFrames = 8
	}
	if tlbEntries <= 0 {
		tlbEntries = 4
	}
	free := make([]int, totalFrames)
	for i := 0; i < totalFrames; i++ {
		free[i] = i
	}
	return &MemoryManager{
		pageSize:    4096,
		totalFrames: totalFrames,
		freeFrames:  free,
		frames:      map[int]frameOwner{},
		pageTables:  map[int]map[uint64]pte{},
		tlb:         make([]tlbEntry, tlbEntries),
	}
}

func (m *MemoryManager) EnsureProcess(pid int) {
	if _, ok := m.pageTables[pid]; !ok {
		m.pageTables[pid] = map[uint64]pte{}
	}
}

func (m *MemoryManager) Protect(pid int, vpn uint64, perm Perm) error {
	pt, ok := m.pageTables[pid]
	if !ok {
		return fmt.Errorf("pid %d has no page table", pid)
	}
	entry, ok := pt[vpn]
	if !ok || !entry.Present {
		return fmt.Errorf("vpn %d not mapped", vpn)
	}
	entry.Perm = perm
	pt[vpn] = entry
	m.invalidateTLB(pid, vpn)
	return nil
}

func (m *MemoryManager) Access(pid int, va uint64, access AccessType) (uint64, string, error) {
	m.EnsureProcess(pid)
	vpn := va / m.pageSize
	offset := va % m.pageSize

	if frame, perm, ok := m.tlbLookup(pid, vpn); ok {
		m.faults.TLBHit++
		if !allowed(perm, access) {
			m.faults.Permission++
			return 0, "permission", fmt.Errorf("permission fault pid=%d vpn=%d access=%s", pid, vpn, access)
		}
		return uint64(frame)*m.pageSize + offset, "", nil
	}
	m.faults.TLBMiss++

	pt := m.pageTables[pid]
	entry, ok := pt[vpn]
	fault := ""
	if !ok || !entry.Present {
		m.faults.NotPresent++
		fault = "not_present"
		frame, err := m.allocateFrame(pid, vpn)
		if err != nil {
			return 0, "", err
		}
		entry = pte{Frame: frame, Perm: Perm{Read: true, Write: true}, Present: true}
		pt[vpn] = entry
	}

	if !allowed(entry.Perm, access) {
		m.faults.Permission++
		return 0, "permission", fmt.Errorf("permission fault pid=%d vpn=%d access=%s", pid, vpn, access)
	}

	m.tlbInsert(pid, vpn, entry.Frame, entry.Perm)
	return uint64(entry.Frame)*m.pageSize + offset, fault, nil
}

func (m *MemoryManager) allocateFrame(pid int, vpn uint64) (int, error) {
	if len(m.freeFrames) > 0 {
		frame := m.freeFrames[0]
		m.freeFrames = m.freeFrames[1:]
		m.frames[frame] = frameOwner{PID: pid, VPN: vpn}
		m.residentQ = append(m.residentQ, frame)
		return frame, nil
	}

	if len(m.residentQ) == 0 {
		return 0, fmt.Errorf("no frames available")
	}

	evictFrame := m.residentQ[0]
	m.residentQ = m.residentQ[1:]
	owner := m.frames[evictFrame]
	ownerPT := m.pageTables[owner.PID]
	old := ownerPT[owner.VPN]
	old.Present = false
	ownerPT[owner.VPN] = old
	m.invalidateTLB(owner.PID, owner.VPN)

	m.frames[evictFrame] = frameOwner{PID: pid, VPN: vpn}
	m.residentQ = append(m.residentQ, evictFrame)
	return evictFrame, nil
}

func (m *MemoryManager) tlbLookup(pid int, vpn uint64) (int, Perm, bool) {
	for _, entry := range m.tlb {
		if entry.Valid && entry.PID == pid && entry.VPN == vpn {
			return entry.Frame, entry.Perm, true
		}
	}
	return 0, Perm{}, false
}

func (m *MemoryManager) tlbInsert(pid int, vpn uint64, frame int, perm Perm) {
	for i := range m.tlb {
		if m.tlb[i].Valid && m.tlb[i].PID == pid && m.tlb[i].VPN == vpn {
			m.tlb[i] = tlbEntry{Valid: true, PID: pid, VPN: vpn, Frame: frame, Perm: perm}
			return
		}
	}
	idx := m.tlbNext % len(m.tlb)
	m.tlb[idx] = tlbEntry{Valid: true, PID: pid, VPN: vpn, Frame: frame, Perm: perm}
	m.tlbNext++
}

func (m *MemoryManager) invalidateTLB(pid int, vpn uint64) {
	for i := range m.tlb {
		if m.tlb[i].Valid && m.tlb[i].PID == pid && m.tlb[i].VPN == vpn {
			m.tlb[i].Valid = false
		}
	}
}

func (m *MemoryManager) Snapshot() MemorySnapshot {
	frames := make([]FrameSnapshot, 0, m.totalFrames)
	for frame := 0; frame < m.totalFrames; frame++ {
		owner, ok := m.frames[frame]
		if !ok {
			frames = append(frames, FrameSnapshot{Frame: frame})
			continue
		}
		frames = append(frames, FrameSnapshot{Frame: frame, PID: owner.PID, VPN: owner.VPN})
	}

	tlb := make([]TLBSnapshot, 0, len(m.tlb))
	for i, entry := range m.tlb {
		if !entry.Valid {
			continue
		}
		tlb = append(tlb, TLBSnapshot{Slot: i, PID: entry.PID, VPN: entry.VPN, Frame: entry.Frame})
	}
	sort.Slice(tlb, func(i, j int) bool { return tlb[i].Slot < tlb[j].Slot })

	return MemorySnapshot{PageSize: m.pageSize, TotalFrames: m.totalFrames, Frames: frames, TLB: tlb, Faults: m.faults}
}

func allowed(perm Perm, access AccessType) bool {
	if access == AccessRead {
		return perm.Read
	}
	if access == AccessWrite {
		return perm.Write
	}
	return false
}
