---
name: repo-workflow
description: Follow repository execution standards for commands, style, architecture boundaries, and minimal-change delivery.
compatibility: opencode
metadata:
  scope: repo-standards
  source: AGENTS.md
---

## What I do
- Apply project command conventions and verification flow.
- Keep changes scoped and aligned with existing naming and structure.
- Update docs when command or workflow behavior changes.

## When to use me
- Use at task start to align implementation approach.
- Use before finalizing to run required checks.

## Repository standards snapshot
- Prefer `make` targets when available.
- Backend core should stay independent from transport/UI.
- Lesson logic observes state and uses command APIs.
- Keep policy modules pluggable (scheduler/replacement).

## Frontend workflow expectations
- When `web/` is present, run frontend checks before handoff.
- Required: `pnpm --dir web eslint .`, `pnpm --dir web tsc --noEmit`, `pnpm --dir web vitest run`, `pnpm --dir web build`.
- Keep UI architecture adapter-driven: transport -> typed DTOs -> selectors -> components.

## Verification checklist
- Formatting and lint completed for touched stack.
- Targeted tests executed first.
- Full relevant suite executed before handoff.
