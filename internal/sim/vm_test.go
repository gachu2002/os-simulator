package sim

import "testing"

func TestMemoryFaultCountsAndFrameOccupancy(t *testing.T) {
	e := NewEngine(77, 0)
	e.ConfigureMemory(2, 2)

	commands := []Command{
		{Name: "spawn", Process: "vm", Program: "ACCESS 0x0 r; ACCESS 0x1000 r; ACCESS 0x2000 r; ACCESS 0x0 r; EXIT"},
		{Name: "step", Count: 12},
	}
	if err := e.ExecuteAll(commands); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	m := e.MemoryView()
	if m.Faults.NotPresent != 4 {
		t.Fatalf("not_present faults=%d want=4", m.Faults.NotPresent)
	}

	if len(m.Frames) != 2 {
		t.Fatalf("frames=%d want=2", len(m.Frames))
	}

	if m.Frames[0].PID != 1 || m.Frames[0].VPN != 2 {
		t.Fatalf("frame0 occupancy=%+v want pid=1 vpn=2", m.Frames[0])
	}
	if m.Frames[1].PID != 1 || m.Frames[1].VPN != 0 {
		t.Fatalf("frame1 occupancy=%+v want pid=1 vpn=0", m.Frames[1])
	}
}

func TestPermissionFaultDeterminism(t *testing.T) {
	mm := NewMemoryManager(2, 2)
	mm.EnsureProcess(1)
	if _, _, err := mm.Access(1, 0, AccessRead); err != nil {
		t.Fatalf("initial access failed: %v", err)
	}
	if err := mm.Protect(1, 0, Perm{Read: true, Write: false}); err != nil {
		t.Fatalf("protect failed: %v", err)
	}
	if _, fault, err := mm.Access(1, 0, AccessWrite); err == nil {
		t.Fatalf("expected permission fault")
	} else if fault != "permission" {
		t.Fatalf("fault=%s want permission", fault)
	}

	snap := mm.Snapshot()
	if snap.Faults.Permission != 1 {
		t.Fatalf("permission faults=%d want=1", snap.Faults.Permission)
	}
}
