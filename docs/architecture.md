# Architecture Snapshot

## Core Packages

- `internal/sim`: deterministic simulator core, schedulers, VM, syscall path, devices/IRQs, filesystem, replay
- `internal/lessons`: staged mission catalog (observe/diagnose/apply), prerequisite gates, validators, hint progression, progress/analytics
- `internal/transport/realtime`: HTTP + WebSocket transport for sessions, lessons, and persisted lesson progress analytics
- `internal/platform/db`: Postgres pool bootstrap (`pgxpool`) used by optional lesson-progress persistence
- `internal/db/sqlc`: generated typed query layer from `sqlc` config and SQL files
- `cmd/simcli`: headless runner for simulation, replay, lesson execution, and analytics
- `cmd/server`: realtime transport server entrypoint for browser sessions
- `web`: React + TypeScript UI with mode routes (`/path`, `/sandbox`, `/challenge`, `/progress`) over deterministic DTO snapshots

## Tooling and Runtime Adapters

- backend HTTP router uses `chi`
- backend server logging uses structured `zap`
- DB workflow uses `golang-migrate` + `sqlc` with PostgreSQL/pgx
- web client HTTP state uses `TanStack Query`
- web styling baseline uses `Tailwind CSS`

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

- trace events emitted for kernel/VM/IO/FS events
- lesson-pack analytics available from CLI
- optional CPU profile and runtime trace capture via CLI flags
