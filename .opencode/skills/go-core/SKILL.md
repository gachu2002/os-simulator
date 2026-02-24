---
name: go-core
description: Build and refactor Go simulator components with small interfaces, clear errors, and package-safe boundaries.
compatibility: opencode
metadata:
  stack: go
  scope: backend
---

## What I do
- Implement Go code in `cmd` and `internal` with minimal, focused changes.
- Keep interfaces behavior-oriented and packages cohesive.
- Apply wrapped errors and strong boundary validation.

## When to use me
- Use when adding or changing backend simulator logic in Go.
- Use when code needs better package structure or interface extraction.

## Guardrails
- Keep imports minimal and goimports-compatible.
- Prefer early returns and short functions.
- Avoid broad refactors unless required by the task.

## Verification
- `gofmt -w .`
- `goimports -w .`
- `go vet ./...`
- `go test ./...`
