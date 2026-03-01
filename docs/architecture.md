# Architecture Snapshot

## Core Packages

- `internal/sim`: deterministic simulator core, schedulers, VM, syscall path, devices/IRQs, filesystem, replay
- `internal/lessons`: lesson catalog, prerequisites, deterministic validators, hint progression, and progress/analytics
- `internal/app/challenges`: challenge application service orchestrating lesson prep/grade with session and attempt stores
- `internal/transport/realtime`: HTTP transport for V3 curriculum and challenge lifecycle (`/curriculum/v3`, `/lessons/{lessonID}/learn/v3`, `/lessons/{lessonID}/challenge/v3`, `/challenges/start|action|submit/v3`, `/challenges/{attemptID}/replay/v3`)
- `internal/platform/db`: isolated infrastructure bootstrap (`pgxpool`), optional at runtime through `DATABASE_URL`; no domain package depends on it
- `internal/db/sqlc`: generated typed query layer from `sqlc` config and SQL files
- `cmd/simcli`: headless runner for simulation, replay, lesson execution, and analytics
- `cmd/server`: realtime transport server entrypoint for browser sessions
- `web`: React + TypeScript UI with Home -> Learn -> Challenge workflow over deterministic DTO snapshots

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
- `internal/app/challenges` orchestrates challenge use cases and depends on domain contracts, not HTTP details.
- `internal/transport/realtime` adapts domain state to HTTP contracts.
- `web` consumes immutable DTO snapshots and never mutates simulator internals.
- `internal/platform/*` is infrastructure-only; it is wired from `cmd/*` entrypoints and must not be imported by domain packages.

## Transport Safety Baseline

- request body limits and unknown-field rejection on write endpoints
- CORS origin controls through allowlist configuration
- request correlation with `X-Request-ID`

## Observability Baseline

- trace events emitted for kernel/VM/IO/FS events
- lesson-pack analytics available from CLI
- optional CPU profile and runtime trace capture via CLI flags
