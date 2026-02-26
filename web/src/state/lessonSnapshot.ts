import type { ChallengeGradeResponse } from "../lib/lessonApi";
import type { SnapshotDTO } from "../lib/types";

export function snapshotFromChallengeGrade(
  result: ChallengeGradeResponse,
): SnapshotDTO {
  return {
    protocol_version: "v1alpha1",
    session_id: `challenge:${result.attempt_id}`,
    tick: result.output.tick,
    trace_hash: result.output.trace_hash,
    trace_length: result.output.trace_length,
    processes: result.output.processes,
    metrics: result.output.metrics,
    memory: result.output.memory,
    last_command: `challenge.grade.${result.lesson_id}.stage.${result.stage_index}`,
  };
}
