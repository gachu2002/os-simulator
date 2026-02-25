package sim

import (
	"fmt"
	"strconv"
	"strings"
)

func (e *Engine) handleEvent(ev Event) {
	if !strings.HasPrefix(ev.Kind, "irq.") {
		return
	}

	requestID, ok := parseRequestID(ev.Data)
	if !ok {
		return
	}

	req, ok := e.devices.Complete(requestID)
	if !ok {
		return
	}

	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "irq.handle", Data: fmt.Sprintf("req=%d pid=%d device=%s", req.ID, req.PID, req.Device)})
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "io.complete", Data: fmt.Sprintf("req=%d pid=%d fd=%d op=%s n=%d", req.ID, req.PID, req.FD, req.Op, req.Bytes)})

	if of, ok := procOpenFile(e.procs, req.PID, req.FD); ok {
		if req.Op == SysRead {
			data, blocks, nextOffset, err := e.fs.ReadInode(of.InodeID, req.Bytes, of.Offset)
			if err == nil {
				of.Offset = nextOffset
				setProcOpenFile(e.procs, req.PID, req.FD, of)
				e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "fs.read", Data: fmt.Sprintf("pid=%d fd=%d bytes=%d", req.PID, req.FD, len(data))})
				e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "fs.blockmap", Data: fmt.Sprintf("pid=%d fd=%d blocks=%v", req.PID, req.FD, blocks)})
			}
		}

		if req.Op == SysWrite {
			payload := []byte(strings.Repeat("w", req.Bytes))
			written, blocks, nextOffset := e.fs.WriteInode(of.InodeID, payload, of.Offset)
			of.Offset = nextOffset
			setProcOpenFile(e.procs, req.PID, req.FD, of)
			e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "fs.write", Data: fmt.Sprintf("pid=%d fd=%d bytes=%d", req.PID, req.FD, written)})
			e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "fs.blockmap", Data: fmt.Sprintf("pid=%d fd=%d blocks=%v", req.PID, req.FD, blocks)})
		}
	}

	proc, ok := e.procs.Get(req.PID)
	if !ok || proc.State != ProcStateBlocked || proc.BlockedOn != "io" {
		return
	}

	_ = proc.transition(ProcStateReady)
	proc.BlockedOn = ""
	proc.BlockedUntil = 0
	e.scheduler.OnReady(proc.PID, true)
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "proc.wakeup", Data: fmt.Sprintf("pid=%d source=irq", proc.PID)})
}

func parseRequestID(data string) (int, bool) {
	if !strings.HasPrefix(data, "req=") {
		return 0, false
	}

	v, err := strconv.Atoi(strings.TrimPrefix(data, "req="))
	if err != nil {
		return 0, false
	}

	return v, true
}

func procOpenFile(pt *ProcessTable, pid, fd int) (OpenFile, bool) {
	proc, ok := pt.Get(pid)
	if !ok {
		return OpenFile{}, false
	}
	of, ok := proc.OpenFiles[fd]
	return of, ok
}

func setProcOpenFile(pt *ProcessTable, pid, fd int, of OpenFile) {
	proc, ok := pt.Get(pid)
	if !ok {
		return
	}
	proc.OpenFiles[fd] = of
}
