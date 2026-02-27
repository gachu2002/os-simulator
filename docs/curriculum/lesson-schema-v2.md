# Lesson Schema v2 (Lesson-Centric)

This schema defines the content model for one OSTEP lesson with:

- one Learn page (theory only)
- one Challenge page (actions, visualization, goal+submit)

This replaces the old stage-centric authoring pattern for curriculum content.

## Design Goals

- Ground theory and challenge directly in OSTEP chapter references.
- Keep challenge grading deterministic and machine-checkable.
- Keep authoring modular and easy to maintain.
- Make failure feedback map to concrete conceptual gaps.

## Top-Level Object

- `schema_version`: string, must be `lesson-v2`
- `lesson`: lesson object

## Lesson Object

- `id`: stable lesson id
- `title`: learner-facing title
- `section_id`: one of `intro`, `cpu-virtualization`, `memory-virtualization`, `concurrency`, `persistence`, `distributed`, `security`, `capstone`
- `difficulty`: `foundation | intermediate | advanced`
- `estimated_minutes`: integer
- `chapter_refs`: array of OSTEP chapter references, e.g. `cpu-sched`, `vm-tlbs`
- `capability_status`: `engine-ready | engine-partial | theory-only`
- `prerequisites`: array of lesson ids
- `learn`: learn content object
- `challenge`: challenge content object

## Learn Content

- `core_idea`: short concept framing
- `mechanism_steps`: ordered array (minimum 4)
- `worked_example`: concrete run explanation tied to expected trace/metrics
- `common_mistakes`: array (minimum 3)
- `pre_challenge_checklist`: array (minimum 3)

## Challenge Content

- `objective`: concise statement of what the learner does
- `goal`: explicit success target in learner language
- `scenario`: deterministic scenario setup
- `actions`: actions panel contract
- `visualization`: visualization panel contract
- `submission`: submit and grading contract
- `hints`: failure-to-hint mapping

### Scenario

- `seed`: deterministic seed
- `base_config`: policy and resource defaults
- `bootstrap_commands`: deterministic setup commands run before learner interaction

### Actions

- `allowed_commands`: array of command ids
- `action_descriptions`: command-to-description map
- `limits`: step/policy/config budgets

### Visualization

- `panels`: ordered panel ids
- `expected_visual_cues`: array of cues the learner should verify

### Submission

- `pass_conditions`: learner-facing checklist
- `validators`: machine validator specs
- `pass_rule`: fixed string: `all_required_and_no_forbidden_and_budget_ok`

### Hints

`hints` is keyed by validator name and has 3 levels:

- `nudge`
- `concept`
- `explicit`

## Validator Spec v2

- `name`: stable validator id
- `type`: one of:
  - `trace_contains_all`
  - `trace_order`
  - `trace_count_eq`
  - `trace_count_lte`
  - `no_event`
  - `metric_eq`
  - `metric_lte`
  - `metric_gte`
  - `fault_eq`
  - `fault_lte`
  - `fs_ok`
  - `budget_ok` (reserved for next engine step)
- `key`: metric/fault key where relevant
- `number`: numeric threshold where relevant
- `values`: event list where relevant
- `required`: boolean (default true)

## Minimal Example (YAML)

```yaml
schema_version: lesson-v2
lesson:
  id: l01-sched-rr-basics
  title: Round Robin Dispatch Basics
  section_id: cpu-virtualization
  difficulty: foundation
  estimated_minutes: 20
  chapter_refs: [cpu-sched]
  capability_status: engine-ready
  prerequisites: []
  learn:
    core_idea: RR rotates CPU ownership with a fixed quantum.
    mechanism_steps:
      - Dispatcher selects next ready process.
      - Process runs until quantum expires or blocks.
      - Timer interrupt triggers preemption on quantum expiry.
      - Process re-enters ready queue and rotation continues.
    worked_example: Two equal CPU-bound jobs alternate dispatch and both finish.
    common_mistakes:
      - Assuming higher context-switch count always means better throughput.
      - Looking only at completion count and ignoring response time.
      - Ignoring trace order when explaining fairness.
    pre_challenge_checklist:
      - I can explain why preemption happens.
      - I can identify dispatch and compute events in trace.
      - I know which metric captures responsiveness.
  challenge:
    objective: Run the workload and demonstrate RR alternation behavior.
    goal: Show alternating dispatch with both processes completing.
    scenario:
      seed: 11
      base_config: {policy: rr, quantum: 2, frames: 8, tlb_entries: 4}
      bootstrap_commands:
        - {name: spawn, process: p1, program: "COMPUTE 4; EXIT"}
        - {name: spawn, process: p2, program: "COMPUTE 4; EXIT"}
    actions:
      allowed_commands: [step, run, pause, reset]
      action_descriptions:
        step: Advance one deterministic tick.
        run: Advance multiple ticks quickly.
        pause: Pause and inspect state.
        reset: Reset to bootstrap state.
      limits: {max_steps: 24, max_policy_changes: 0, max_config_changes: 0}
    visualization:
      panels: [timeline, process_queues, metrics]
      expected_visual_cues:
        - Dispatch alternates between p1 and p2 before exits.
        - Completed process count reaches 2.
    submission:
      pass_rule: all_required_and_no_forbidden_and_budget_ok
      pass_conditions:
        - Trace contains dispatch and compute events for both jobs.
        - Both processes complete in budget.
      validators:
        - {name: trace_core, type: trace_contains_all, values: [proc.dispatch, proc.compute], required: true}
        - {name: completed, type: metric_eq, key: completed_processes, number: 2, required: true}
    hints:
      trace_core:
        nudge: Check event kinds in timeline first.
        concept: RR fairness is visible in dispatch cadence.
        explicit: Step through and verify alternating dispatch before submit.
      completed:
        nudge: Verify both jobs receive enough CPU time.
        concept: Completion is required even when response improves.
        explicit: Run additional ticks until completed_processes reaches 2.
```

## Migration Notes

- Existing runtime may continue to use internal `Stage` structs during migration.
- Authoring source of truth should move to this lesson-centric schema.
- UI must render Learn and Challenge from this model without stage-tab logic.
