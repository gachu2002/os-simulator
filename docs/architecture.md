# Architecture Snapshot

## Core Packages

- `internal/sim`: deterministic simulator core, schedulers, VM, syscall path, devices/IRQs, filesystem, replay
- `internal/lessons`: lesson DSL/catalog, validators, hint progression, progress/analytics
- `internal/transport/realtime`: HTTP + WebSocket session transport, command validation, immutable snapshot DTO stream
- `cmd/simcli`: headless runner for simulation, replay, lesson execution, and analytics
- `cmd/server`: realtime transport server entrypoint for browser sessions
- `web`: React + TypeScript control UI for session creation, run/pause/step/reset, and event log/status panels

## Determinism Contract

- event ordering is stable by `(tick, sequence)`
- no wall-clock dependence in simulation behavior
- deterministic replay hash validated by golden tests
- deterministic regression suite enforced in CI

## Boundary Contract

- `internal/sim` and `internal/lessons` are domain core; they do not depend on transport or UI.
- `internal/transport/realtime` adapts domain state to HTTP/WS contracts.
- `web` consumes immutable DTO snapshots and never mutates simulator internals.

## Transport Safety Baseline

- request body limits and unknown-field rejection on write endpoints
- CORS/WS origin controls through allowlist configuration
- request correlation with `X-Request-ID`

## Observability Baseline

- trace events emitted for kernel/VM/IO/FS milestones
- lesson-pack analytics available from CLI
- optional CPU profile and runtime trace capture via CLI flags
