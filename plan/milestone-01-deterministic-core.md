# Milestone 01: Deterministic Core

## Goal

Build a reproducible simulator engine loop.

## Scope

- Simulation clock
- Event queue
- Snapshot manager
- Replay log

## Deliverables

- CLI runner for headless simulation steps
- Golden trace tests for replay consistency

## Exit Criteria

- Same seed + same commands produce identical trace hash

## Key Risks

- Hidden nondeterminism from map iteration
- Hidden nondeterminism from wall-clock usage

## Suggested First Tasks

1. Define deterministic event ordering contract `(tick, sequence)`.
2. Implement event queue and simulation clock.
3. Add replay writer/reader and snapshot checkpoints.
4. Create golden trace test harness in Go.
