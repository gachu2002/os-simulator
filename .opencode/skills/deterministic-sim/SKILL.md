---
name: deterministic-sim
description: Enforce deterministic simulator behavior through stable event ordering, seeded randomness, and replay-safe state transitions.
compatibility: opencode
metadata:
  stack: go
  focus: determinism
---

## What I do
- Design tick-based event handling with stable `(tick, sequence)` ordering.
- Remove wall-clock and non-seeded randomness from core logic.
- Keep snapshots and traces reproducible across runs.

## When to use me
- Use for scheduler, VM, interrupt, and event queue changes.
- Use when a test flakes or replay output diverges.

## Rules
- No correctness-critical reliance on map iteration order.
- Deterministic tie-breakers are mandatory.
- State transitions must be explicit and testable.

## Verification
- Add or update replay/trace tests.
- `go test ./...`
- Re-run the same test twice and confirm identical results.
