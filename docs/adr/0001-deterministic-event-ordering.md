# ADR 0001: Deterministic Event Ordering

## Status

Accepted

## Decision

Simulator events are ordered by stable `(tick, sequence)` semantics. Replay hash is used as determinism regression signal.

## Rationale

Educational and testing value depends on reproducible traces across runs and environments.

## Consequences

- No wall-clock dependence in core domain behavior.
- Determinism-sensitive changes require replay/hash tests.
