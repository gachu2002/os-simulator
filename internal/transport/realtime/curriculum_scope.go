package realtime

import (
	"fmt"

	contentv3 "os-simulator-plan/internal/content/v3"
)

const sectionVirtualizationCPU = "virtualization-cpu"

var activeCPULessonOrder = []string{
	"l01-process-basics",
	"l02-process-api-fork-exec-wait",
	"l03-limited-direct-execution",
	"l04-cpu-scheduling-basics",
	"l05-round-robin",
	"l06-mlfq",
	"l07-lottery-stride",
	"l08-multi-cpu-scheduling",
}

var activeCPULessonSet = map[string]struct{}{
	"l01-process-basics":             {},
	"l02-process-api-fork-exec-wait": {},
	"l03-limited-direct-execution":   {},
	"l04-cpu-scheduling-basics":      {},
	"l05-round-robin":                {},
	"l06-mlfq":                       {},
	"l07-lottery-stride":             {},
	"l08-multi-cpu-scheduling":       {},
}

var v3ToEngineLessonID = map[string]string{
	"l01-process-basics":             "l01-sched-rr-basics",
	"l02-process-api-fork-exec-wait": "l02-sched-fifo-baseline",
	"l03-limited-direct-execution":   "l03-sched-mlfq-balance",
	"l04-cpu-scheduling-basics":      "l04-response-under-rr",
	"l05-round-robin":                "l05-throughput-shared-cpu",
	"l06-mlfq":                       "l06-preemption-check",
	"l07-lottery-stride":             "l06b-lottery-tradeoffs",
	"l08-multi-cpu-scheduling":       "l06c-quantum-response-tuning",
}

var engineToV3LessonID = map[string]string{
	"l01-sched-rr-basics":          "l01-process-basics",
	"l02-sched-fifo-baseline":      "l02-process-api-fork-exec-wait",
	"l03-sched-mlfq-balance":       "l03-limited-direct-execution",
	"l04-response-under-rr":        "l04-cpu-scheduling-basics",
	"l05-throughput-shared-cpu":    "l05-round-robin",
	"l06-preemption-check":         "l06-mlfq",
	"l06b-lottery-tradeoffs":       "l07-lottery-stride",
	"l06c-quantum-response-tuning": "l08-multi-cpu-scheduling",
}

func isActiveCPULesson(id string) bool {
	_, ok := activeCPULessonSet[id]
	return ok
}

func resolveEngineLessonID(v3LessonID string) (string, bool) {
	engineID, ok := v3ToEngineLessonID[v3LessonID]
	return engineID, ok
}

func resolveV3LessonID(engineLessonID string) (string, bool) {
	v3ID, ok := engineToV3LessonID[engineLessonID]
	return v3ID, ok
}

func validateCPUCurriculumScope(curriculum contentv3.Curriculum) error {
	if len(curriculum.Sections) != 1 {
		return fmt.Errorf("expected one section, got %d", len(curriculum.Sections))
	}
	section := curriculum.Sections[0]
	if section.ID != sectionVirtualizationCPU {
		return fmt.Errorf("section id=%q want=%q", section.ID, sectionVirtualizationCPU)
	}
	if len(section.Lessons) != len(activeCPULessonOrder) {
		return fmt.Errorf("lesson count=%d want=%d", len(section.Lessons), len(activeCPULessonOrder))
	}
	for idx, lesson := range section.Lessons {
		if lesson.ID != activeCPULessonOrder[idx] {
			return fmt.Errorf("lesson[%d]=%q want=%q", idx, lesson.ID, activeCPULessonOrder[idx])
		}
	}
	return nil
}
