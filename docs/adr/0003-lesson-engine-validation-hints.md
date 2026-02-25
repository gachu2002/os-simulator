# ADR 0003: Lesson Validation and Hint Progression

## Status

Accepted

## Decision

Lesson stages are validated by deterministic output predicates. Failed attempts progress hint level from nudge to concept to explicit.

## Rationale

This keeps pedagogical feedback predictable, testable, and aligned with deterministic replay behavior.

## Consequences

- Validators and hint levels are part of API-visible lesson outcomes.
- Lesson tests must cover pass/fail and hint progression behavior.
