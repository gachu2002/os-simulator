# OS Simulator Plan (v1.0 RC)

Deterministic OSTEP-aligned simulator core implemented in Go with:

- deterministic engine + replay log + golden hash testing
- process lifecycle and schedulers (FIFO, RR, MLFQ)
- syscall/trap path with async device interrupts
- virtual memory (VA->PA, TLB, faults, FIFO replacement)
- filesystem path traversal + block mapping
- lesson engine + 20-lesson OSTEP pack with analytics

## Quick Start

```bash
go test ./...
go run ./cmd/simcli -program "SYSCALL open /docs/readme.txt; SYSCALL read 4; SYSCALL write 3; SYSCALL exit" -steps 16
go run ./cmd/simcli -run-lesson-pack
```

## Stable Engineering Workflow

Use `make` targets:

- `make fmt` - format Go code
- `make test` - full tests
- `make test-deterministic` - deterministic regression suite
- `make lesson-pack` - lesson-pack analytics smoke
- `make ci` - full CI-equivalent local run
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
