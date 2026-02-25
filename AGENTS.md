# AGENTS.md

## Rule Files Discovery

Checked in repo:

- `.cursorrules`: not found
- `.cursor/rules/`: not found
- `.github/copilot-instructions.md`: not found
  If these appear later, they override this file where instructions conflict.

## Project Context

- Product: interactive OSTEP-aligned OS simulator
- Backend: Go
- Frontend: React + TypeScript
- Critical property: deterministic simulation and replayability

## Tooling Baseline

- Backend routing/logging: `chi` + `zap`
- Backend data path scaffold: PostgreSQL + `pgx` + `sqlc` + `golang-migrate`
- Backend dev reload: `air`
- Frontend stack: `Vite` + `Tailwind CSS` + `shadcn/ui` scaffold + `TanStack Query`
- Frontend quality: `ESLint` + `Prettier` + `Vitest`

## Command Policy

- Prefer `make` targets if available.
- If no `make`, use raw commands below.
- Run targeted tests first, then full suite.

## Build Commands

### Go

- Build all: `go build ./...`
- Build package: `go build ./internal/sim`
- Build server entry (if present): `go build ./cmd/server`
- Build CLI entry (if present): `go build ./cmd/simcli`
- SQLC generate: `make sqlc-generate`
- Migration up: `make db-up`
- Migration down: `make db-down`
- Migration status: `make db-status`
- Migration create: `make db-create name=add_table`
- Dev reload server: `make dev-server`

### Web

- Install deps: `pnpm --dir=web install`
- Dev server: `pnpm --dir=web run dev`
- Build: `pnpm --dir=web run build`
- Preview: `pnpm --dir=web run preview`
- Add shadcn component: `make web-shadcn-add name=button`

## Lint and Format

### Go

- Format: `gofmt -w .`
- Imports: `goimports -w .`
- Vet: `go vet ./...`
- Lint: `golangci-lint run`

### TypeScript / React

- Lint: `pnpm --dir=web run lint`
- Prettier check: `pnpm --dir=web exec prettier --check .`
- Typecheck: `pnpm --dir=web run typecheck`
- Prettier write: `pnpm --dir=web exec prettier --write .`

## Test Commands (Single-Test Focus)

### Go tests

- All tests: `go test ./...`
- Verbose: `go test -v ./...`
- Race: `go test -race ./...`
- One package: `go test ./internal/sim`
- One test: `go test ./internal/sim -run '^TestGoldenTraceHash$'`
- One subtest: `go test ./internal/sim -run 'TestKnownWorkloadMetrics_FIFOvsRR/rr'`
- Benchmark: `go test ./internal/sim -bench . -benchmem`
- Generic single-test form: `go test ./path/to/pkg -run '^TestName$'`

### Frontend unit tests

- All tests: `pnpm --dir=web run test`
- Run once: `pnpm --dir=web exec vitest run`
- One file: `pnpm --dir=web exec vitest run src/components/Timeline.test.tsx`
- One test title: `pnpm --dir=web exec vitest run -t "renders ready queue"`

### E2E tests

- All specs: `pnpm --dir=web exec playwright test`
- One spec: `pnpm --dir=web exec playwright test tests/scheduler.spec.ts`
- One test title: `pnpm --dir=web exec playwright test -g "step mode updates Gantt"`

## Verification Sequence

### Backend-only

1. `gofmt -w .` and `goimports -w .`
2. `golangci-lint run`
3. Run targeted Go test(s)
4. `go test ./...`
5. `go test -race ./...`

### Frontend-only

1. `pnpm --dir=web exec prettier --write .`
2. `pnpm --dir=web run lint`
3. `pnpm --dir=web run typecheck`
4. Run targeted frontend test(s)
5. `pnpm --dir=web run test`

### Cross-cutting

1. Run targeted backend + frontend tests first
2. Run full backend tests
3. Run frontend lint, typecheck, test, build
4. Run `govulncheck ./...` and dependency audits for release-critical changes

## Code Style Guidelines

### Imports

- Keep imports minimal and used.
- Go grouping: standard library, third-party, internal.
- Enforce with `goimports`.
- In TS, prefer named imports; avoid wildcard imports unless justified.

### Formatting

- Do not hand-format against tooling output.
- Go uses `gofmt` + `goimports`.
- TS/JS/MD use Prettier.
- Prefer early returns over deep nesting.

### Types and Interfaces

- Use explicit types at boundaries.
- In TS, avoid `any`; use `unknown` then narrow.
- In Go, define small behavior-focused interfaces.
- Treat snapshots/DTOs as immutable after emission.
- Model state transitions explicitly.

### Naming

- Go exported: `PascalCase`; internal: `camelCase`.
- TS components/types: `PascalCase`.
- TS variables/functions/hooks: `camelCase`.
- Constants: `UPPER_SNAKE_CASE`.
- Keep acronyms consistent: `PID`, `TLB`, `IRQ`, `CPU`.

### Error Handling

- Never swallow errors silently.
- Wrap Go errors with context: `fmt.Errorf("context: %w", err)`.
- Use sentinel/typed errors only for caller branching.
- Validate input at transport boundaries (HTTP/WS).
- Return actionable user-facing error messages.

### Concurrency and Determinism

- No wall-clock dependence in simulator core.
- No unseeded randomness in domain logic.
- Keep stable ordering for equal-priority events.
- Do not rely on map iteration order in correctness-critical paths.
- Add replay/trace tests for determinism-sensitive changes.

## Architecture Boundaries

- Simulator core stays independent from transport/UI.
- UI consumes snapshots/DTOs and does not mutate core internals.
- Lesson logic observes state and acts through command APIs.
- Keep policy modules pluggable (scheduler, replacement).

## Agent Expectations

- Make the smallest coherent change that solves the task.
- Preserve existing patterns and names.
- Avoid broad refactors without direct requirement.
- Update docs when commands/workflows/structure change.
- If repo reality diverges from this guide, update `AGENTS.md` in the same change.

## Generated Artifacts Policy

- Do not hand-edit generated files when a framework/tool is the source of truth.
- For SQL access layer changes: edit SQL in `db/query`/`db/schema`, then run `make sqlc-generate`.
- For migrations: create new files via `make db-create name=<migration_name>`.
- For shadcn/ui components: generate via `make web-shadcn-add name=<component>` (uses `shadcn@latest`).
- If generated output changes are expected, include generated files in the same change.
