# Engineering Workflow

This repository follows a small-team workflow designed to scale to mid-size orgs without heavy process overhead.

## Branch and PR Model

- Protect `main`; all changes land through pull requests.
- Keep PRs scoped to one concern (feature, bugfix, refactor, or docs).
- Prefer short-lived branches and merge within 1-2 days.
- Use conventional commit prefixes (`feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`).

## Definition of Done

A change is done when all of the following are true:

- architecture boundaries remain intact (`internal/sim` and `internal/lessons` do not depend on transport or UI)
- targeted tests pass for touched behavior
- full relevant suites pass (`make ci-go` and/or `make ci-web`)
- docs are updated for contract/workflow/command changes
- no accidental dead code or generated-file drift is introduced

## Review Expectations

- include problem statement, approach, and risk in PR description
- include validation commands and outcomes
- call out behavior changes in API/DTO/contracts explicitly
- request CODEOWNERS review for owned paths

## Refactor Policy

- refactor in small slices with behavior-preserving tests
- keep each refactor PR revertable
- introduce adapter seams before moving domain logic
- avoid broad package moves without ADR updates

## Redundant and Legacy Code Cleanup

Use this sequence before deleting code/docs:

1. Prove unused with repository search and test coverage.
2. Remove only one coherent legacy area per PR.
3. Run targeted checks first, then full affected suites.
4. Document removals in PR notes with rollback instructions.

Suggested command gate:

- `make audit-unused`

## Infrastructure Isolation

- `internal/platform/*` is infrastructure-only and wired from `cmd/*`.
- Domain and application packages must not import infrastructure packages.
- Optional scaffolding (DB, generated access layer) is allowed when isolated and documented.

## Release Gate

- Follow `docs/release-checklist.md` for candidate and release preparation.
- CI (`.github/workflows/ci.yml`) and deploy smoke (`.github/workflows/deploy-smoke.yml`) must be green for production promotion.
