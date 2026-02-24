---
name: ostep-kernel
description: Implement OSTEP-aligned kernel behavior for process state, scheduling, syscalls, memory, and interrupt-driven wakeups.
compatibility: opencode
metadata:
  stack: go
  focus: os-model
---

## What I do
- Build process lifecycle, scheduler policies, syscall flow, and wakeup paths.
- Keep kernel behavior aligned with educational scenarios.
- Expose clean DTO-ready outputs without leaking core internals.

## When to use me
- Use when implementing FIFO/RR/MLFQ, traps, page faults, or FS basics.
- Use when lesson goals depend on kernel state transitions.

## Modeling expectations
- Preserve `new -> ready -> running -> blocked -> ready -> terminated` style transitions.
- Keep policy modules pluggable.
- Validate syscall inputs at system boundaries.

## Verification
- Targeted package tests first:
  - `go test ./internal/kernel/sched -run '^TestRoundRobinBasic$'`
  - `go test ./internal/kernel/sched -run 'TestRoundRobinBasic/quantum_4'`
- Then full suite: `go test ./...`
