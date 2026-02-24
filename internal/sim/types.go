package sim

type Tick uint64

type Event struct {
	Tick     Tick   `json:"tick"`
	Sequence uint64 `json:"sequence"`
	Kind     string `json:"kind"`
	Data     string `json:"data,omitempty"`
}

type TraceEvent = Event

type Command struct {
	Name    string `json:"name"`
	Count   int    `json:"count,omitempty"`
	Tick    Tick   `json:"tick,omitempty"`
	Kind    string `json:"kind,omitempty"`
	Data    string `json:"data,omitempty"`
	Process string `json:"process,omitempty"`
	Program string `json:"program,omitempty"`
	Policy  string `json:"policy,omitempty"`
	Quantum int    `json:"quantum,omitempty"`
}

type ProcessSnapshot struct {
	PID          int       `json:"pid"`
	Name         string    `json:"name"`
	State        ProcState `json:"state"`
	PC           int       `json:"pc"`
	BlockedUntil Tick      `json:"blocked_until,omitempty"`
}

type GanttSlice struct {
	Tick Tick `json:"tick"`
	PID  int  `json:"pid"`
}

type ProcessMetric struct {
	PID          int    `json:"pid"`
	Name         string `json:"name"`
	ResponseTime Tick   `json:"response_time"`
	Turnaround   Tick   `json:"turnaround"`
	RunTicks     Tick   `json:"run_ticks"`
	WaitTicks    Tick   `json:"wait_ticks"`
}

type SchedulingMetrics struct {
	Policy               string          `json:"policy"`
	Quantum              int             `json:"quantum,omitempty"`
	TotalTicks           Tick            `json:"total_ticks"`
	CompletedProcesses   int             `json:"completed_processes"`
	AvgResponseTime      float64         `json:"avg_response_time"`
	AvgTurnaroundTime    float64         `json:"avg_turnaround_time"`
	ThroughputPer100Tick float64         `json:"throughput_per_100_ticks"`
	FairnessJainIndex    float64         `json:"fairness_jain_index"`
	Processes            []ProcessMetric `json:"processes"`
	Gantt                []GanttSlice    `json:"gantt"`
}

type FrameSnapshot struct {
	Frame int    `json:"frame"`
	PID   int    `json:"pid,omitempty"`
	VPN   uint64 `json:"vpn"`
}

type TLBSnapshot struct {
	Slot  int    `json:"slot"`
	PID   int    `json:"pid"`
	VPN   uint64 `json:"vpn"`
	Frame int    `json:"frame"`
}

type FaultCounters struct {
	NotPresent int `json:"not_present"`
	Permission int `json:"permission"`
	TLBHit     int `json:"tlb_hit"`
	TLBMiss    int `json:"tlb_miss"`
}

type MemorySnapshot struct {
	PageSize    uint64          `json:"page_size"`
	TotalFrames int             `json:"total_frames"`
	Frames      []FrameSnapshot `json:"frames"`
	TLB         []TLBSnapshot   `json:"tlb"`
	Faults      FaultCounters   `json:"faults"`
}
