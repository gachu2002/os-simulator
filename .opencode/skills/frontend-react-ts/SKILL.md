---
name: frontend-react-ts
description: Build and maintain typed React UI for simulator controls and timeline views without mutating simulator internals.
compatibility: opencode
metadata:
  stack: react-typescript
  scope: frontend
---

## What I do
- Implement React + TypeScript UI for run/pause/step/reset and visual state panels.
- Keep strict typing for event payloads and snapshot DTOs.
- Preserve separation between UI state and simulator core domain.

## When to use me
- Use when adding or updating components, hooks, or event-stream-driven screens.
- Use when wiring frontend views to backend simulation snapshots.

## Guardrails
- Avoid `any`; use `unknown` and narrow.
- Prefer named imports and consistent component naming.
- Do not mutate data that represents immutable snapshots.

## Verification
- `pnpm --dir web prettier --write .`
- `pnpm --dir web eslint .`
- `pnpm --dir web tsc --noEmit`
- `pnpm --dir web vitest run`
