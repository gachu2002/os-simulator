# AGENTS.md

## Purpose
Default working guide for agentic coding in this repository.
Current state is planning-heavy: `IMPLEMENTATION_PLAN.md`, `SKILLS.md`, `plan/`, and `.opencode/skills/`.

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
- Architecture reference: `IMPLEMENTATION_PLAN.md`

## Command Policy
- Prefer `make` targets if available.
- If no `make`, use raw commands below.
- Run targeted tests first, then full suite.

## Build Commands

### Go
- Build all: `go build ./...`
- Build package: `go build ./internal/sim/...`
- Build server entry (if present): `go build ./cmd/server`
- Build CLI entry (if present): `go build ./cmd/simcli`

### Web
- Install deps: `pnpm --dir web install`
- Dev server: `pnpm --dir web dev`
- Build: `pnpm --dir web build`
- Preview: `pnpm --dir web preview`

## Lint and Format

### Go
- Format: `gofmt -w .`
- Imports: `goimports -w .`
- Vet: `go vet ./...`
- Lint: `golangci-lint run`

### TypeScript / React
- Lint: `pnpm --dir web eslint .`
- Typecheck: `pnpm --dir web tsc --noEmit`
- Prettier check: `pnpm --dir web prettier --check .`
- Prettier write: `pnpm --dir web prettier --write .`

## Test Commands (Single-Test Focus)

### Go tests
- All tests: `go test ./...`
- Verbose: `go test -v ./...`
- Race: `go test -race ./...`
- One package: `go test ./internal/kernel/sched`
- One test: `go test ./internal/kernel/sched -run '^TestRoundRobinBasic$'`
- One subtest: `go test ./internal/kernel/sched -run 'TestRoundRobinBasic/quantum_4'`
- Benchmark: `go test ./internal/sim -bench . -benchmem`
- Generic single-test form: `go test ./path/to/pkg -run '^TestName$'`

### Frontend unit tests
- All tests: `pnpm --dir web test`
- Run once: `pnpm --dir web vitest run`
- One file: `pnpm --dir web vitest run src/components/Timeline.test.tsx`
- One test title: `pnpm --dir web vitest run -t "renders ready queue"`

### E2E tests
- All specs: `pnpm --dir web playwright test`
- One spec: `pnpm --dir web playwright test tests/scheduler.spec.ts`
- One test title: `pnpm --dir web playwright test -g "step mode updates Gantt"`

## Verification Sequence

### Backend-only
1. `gofmt -w .` and `goimports -w .`
2. `golangci-lint run`
3. Run targeted Go test(s)
4. `go test ./...`

### Frontend-only
1. `pnpm --dir web prettier --write .`
2. `pnpm --dir web eslint .`
3. `pnpm --dir web tsc --noEmit`
4. Run targeted frontend test(s)
5. `pnpm --dir web test`

### Cross-cutting
1. Run targeted backend + frontend tests first
2. Run full backend tests
3. Run frontend lint, typecheck, test, build

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
