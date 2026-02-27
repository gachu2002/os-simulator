# Failure Taxonomy v2

This taxonomy standardizes how validator failures map to learner feedback.

Goal: every failed submission gives a concept-specific correction path.

## Failure Classes

- `trace-missing`: required trace events did not appear.
- `trace-order`: required events appeared in wrong order.
- `trace-frequency`: event count/frequency mismatch.
- `forbidden-event`: an event that should never appear was observed.
- `metric-threshold`: metric eq/lte/gte condition failed.
- `memory-fault`: fault_eq/fault_lte mismatch.
- `filesystem-integrity`: `fs_ok` failed.
- `budget-violation`: exceeded step/policy/config limits (engine support in progress).

## Validator -> Failure Class Mapping

- `trace_contains_all` -> `trace-missing`
- `trace_order` -> `trace-order`
- `trace_count_eq`, `trace_count_lte` -> `trace-frequency`
- `no_event` -> `forbidden-event`
- `metric_eq`, `metric_lte`, `metric_gte` -> `metric-threshold`
- `fault_eq`, `fault_lte` -> `memory-fault`
- `fs_ok` -> `filesystem-integrity`
- `budget_ok` -> `budget-violation`

## Hint Policy

Each failure class has three escalating levels:

- `L1 Nudge`: where to inspect (trace/metrics/memory/fs panel)
- `L2 Concept`: why the model predicts a different outcome
- `L3 Explicit`: concrete correction action

## Runtime Behavior

- Grade engine reports first failed validator as feedback key source.
- If validator-specific hints are configured, use them first.
- Fallback to stage-level hints only when validator-specific hints are absent.

## Authoring Rules

- Every validator in a lesson challenge should have a corresponding hint triple.
- Hint text should reference one mechanism step from Learn page.
- Explicit hints should suggest one change at a time (avoid multi-action instructions).

## Iteration Metrics

Use failure telemetry to improve lesson quality:

- first failing validator frequency
- attempts to pass per lesson
- proportion reaching hint level 3
- repeated failures on same validator after explicit hint

Target thresholds for mature lessons:

- <= 35% of attempts require level-3 hints
- <= 2 median attempts to pass
- <= 10% repeated failure after explicit hint
