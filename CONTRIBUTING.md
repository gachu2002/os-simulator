# Contributing

## Workflow

1. Create a branch from `main`.
2. Keep changes small and focused.
3. Run local verification before opening PR.
4. Open PR with summary, rationale, and validation evidence.

## Local Verification

```bash
make ci
```

Optional local guardrail:

```bash
pre-commit install
```

If only one stack is touched, run targeted checks first, then full relevant suite.

## Commit Style

Use concise conventional-style prefixes:

- `feat:` new capability
- `fix:` bug fix
- `refactor:` structure change without behavior change
- `test:` tests only
- `docs:` documentation only
- `chore:` tooling/workflow upkeep

## Determinism Rules

- No wall-clock behavior in simulator core.
- Stable ordering for equal-priority events.
- Add replay/hash regression tests for determinism-sensitive updates.

## Code Review Expectations

- API boundaries remain explicit and typed.
- Error handling remains actionable and consistent.
- Tests cover happy path + failure path for touched behavior.
