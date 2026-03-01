# OS Simulator Plan (v1.0 RC)

Deterministic OSTEP-aligned simulator core implemented in Go with:

- deterministic engine + replay log + golden hash testing
- process lifecycle and schedulers (FIFO, RR, MLFQ)
- syscall/trap path with async device interrupts
- virtual memory (VA->PA, TLB, faults, FIFO replacement)
- filesystem path traversal + block mapping
- lesson engine with chapter-grounded lesson theory and deterministic challenge grading
- prerequisite-gated curriculum path (Virtualization -> Concurrency -> Persistence)
- basic completion analytics per challenge attempt

## Quick Start

```bash
go test ./...
go run ./cmd/simcli -program "SYSCALL open /docs/readme.txt; SYSCALL read 4; SYSCALL write 3; SYSCALL exit" -steps 16
go run ./cmd/simcli -run-lesson-pack
```

## Realtime Server + Web UI

```bash
go run ./cmd/server -addr :8080
pnpm --dir=web install
pnpm --dir=web run dev
```

Optional infrastructure bootstrap:

- set `DATABASE_URL` to enable Postgres pool bootstrap on server startup
- leave `DATABASE_URL` unset to run in simulator-only mode

Web routes:

- `/`: OSTEP section overview with lesson titles (Introduction/Security shown as coming soon)
- `/lesson/:lessonId/learn`: theory-only lesson page
- `/lesson/:lessonId/challenge`: challenge page with actions, visualization, and goal+submit

## Tooling Baseline

- Backend: `chi`, `pgx` + `sqlc`, `golang-migrate`, `zap`, `golangci-lint`, `air`
- Frontend: `Vite`, `Tailwind CSS`, `TanStack Query`, `ESLint` + `Prettier`, `Vitest`

## Stable Engineering Workflow

Use `make` targets:

- `make fmt` - format Go code
- `make fmt-check` - verify Go formatting
- `make lint` - run Go lint checks
- `make test` - full tests
- `make test-race` - race detector tests
- `make test-coverage` - enforce package coverage targets
- `make test-deterministic` - deterministic regression suite
- `make lesson-pack` - lesson-pack analytics smoke
- `make sqlc-generate` - generate typed DB access code from SQL
- `make sqlc-verify` - verify generated DB code is up to date
- `make ci-go` - Go CI-equivalent local run
- `make ci-web` - web CI-equivalent local run
- `make ci-security` - security CI-equivalent local run
- `make audit-unused` - detect unused/dead code signals (Go + TypeScript)
- `make db-up` / `make db-down` / `make db-status` - run local DB migrations
- `make db-create name=add_feature` - create migration pair
- `make dev-server` - run backend with live reload via air
- `make web-format-check` - run Prettier check for web sources
- `make ci` - full CI-equivalent local run
- `make security` - vulnerability and dependency audit checks
- `make release-check` - CI checks + full build

## Observability and Profiling

`simcli` supports:

- structured observability output: `-emit-observability`
- CPU profile: `-cpu-profile cpu.pprof`
- runtime trace: `-trace-file runtime.trace`

Example:

```bash
go run ./cmd/simcli -program "ACCESS 0x0 r; ACCESS 0x1000 r; EXIT" -steps 12 -emit-observability -cpu-profile cpu.pprof -trace-file runtime.trace
```

## Release Process

See `docs/release-checklist.md`.

## API and Engineering Docs

- API reference: `docs/api.md`
- Architecture: `docs/architecture.md`
- Engineering workflow: `docs/engineering-workflow.md`
- ADR index: `docs/adr/README.md`
- Contribution guide: `CONTRIBUTING.md`
- Security policy: `SECURITY.md`

## Free Deploy + CD

See `docs/deployment-free.md` for the recommended free setup (Render + Cloudflare Pages + GitHub deploy smoke checks).
