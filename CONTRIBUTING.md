# Contributing

## Workflow

1. Create a branch from `main`.
2. Keep changes small and focused.
3. Run local verification before opening PR.
4. Open PR with summary, rationale, and validation evidence.

See `docs/engineering-workflow.md` for repository operating model, review expectations, and cleanup policy.

## Local Verification

```bash
make ci
```

Optional local guardrail:

```bash
pre-commit install
```

If only one stack is touched, run targeted checks first, then full relevant suite.

### Backend pull requests

Run this sequence for backend-only changes:

1. Targeted tests for touched package(s).
2. `make ci-go`
3. If determinism-sensitive behavior changed, rerun `make test-deterministic` once more to confirm stable results.

### Toolchain pinning

Repository `make` targets pin CLI tool versions for reproducibility. When bumping a tool version, update it in `Makefile` and include the reason in the PR.

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

## Generated Artifacts

- Do not hand-edit generated files.
- For SQL access layer changes, edit SQL under `db/` and run `make sqlc-generate`.
- For new DB migrations, run `make db-create name=<migration_name>`.
- Include generated output in the same PR when source changes require regeneration.
