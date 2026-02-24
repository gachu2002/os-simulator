package sim

import "fmt"

const (
	DeviceDisk     = "disk"
	DeviceTerminal = "terminal"
)

type IORequest struct {
	ID         int
	PID        int
	FD         int
	Device     string
	Op         string
	Bytes      int
	CompleteAt Tick
}

type DeviceManager struct {
	diskLatency     Tick
	terminalLatency Tick
	nextRequestID   int
	pending         map[int]IORequest
}

func NewDeviceManager(diskLatency, terminalLatency Tick) *DeviceManager {
	if diskLatency == 0 {
		diskLatency = 3
	}
	if terminalLatency == 0 {
		terminalLatency = 1
	}
	return &DeviceManager{
		diskLatency:     diskLatency,
		terminalLatency: terminalLatency,
		nextRequestID:   1,
		pending:         map[int]IORequest{},
	}
}

func (m *DeviceManager) Submit(now Tick, pid int, fd int, device, op string, bytes int) IORequest {
	latency := m.terminalLatency
	if device == DeviceDisk {
		latency = m.diskLatency
	}
	req := IORequest{
		ID:         m.nextRequestID,
		PID:        pid,
		FD:         fd,
		Device:     device,
		Op:         op,
		Bytes:      bytes,
		CompleteAt: now + latency,
	}
	m.nextRequestID++
	m.pending[req.ID] = req
	return req
}

func (m *DeviceManager) Complete(requestID int) (IORequest, bool) {
	req, ok := m.pending[requestID]
	if !ok {
		return IORequest{}, false
	}
	delete(m.pending, requestID)
	return req, true
}

func IRQEventKind(device string) string {
	return fmt.Sprintf("irq.%s.complete", device)
}
