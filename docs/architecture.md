# Architecture Snapshot

## Core Packages

- `internal/sim`: deterministic simulator core, schedulers, VM, syscall path, devices/IRQs, filesystem, replay
- `internal/lessons`: lesson DSL/catalog, validators, hint progression, progress/analytics
- `cmd/simcli`: headless runner for simulation, replay, lesson execution, and analytics

## Determinism Contract

- event ordering is stable by `(tick, sequence)`
- no wall-clock dependence in simulation behavior
- deterministic replay hash validated by golden tests
- deterministic regression suite enforced in CI

## Observability Baseline

- trace events emitted for kernel/VM/IO/FS milestones
- lesson-pack analytics available from CLI
- optional CPU profile and runtime trace capture via CLI flags
