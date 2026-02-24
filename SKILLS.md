# Skills Index

This repository uses `.opencode/skills/` for task-specific execution guidance.

## Available Skills

- `repo-workflow`: repository command conventions, verification flow, and boundary rules.
- `testing-discipline`: targeted deterministic testing first, then full-suite verification.
- `go-core`: focused backend implementation for Go simulator packages.
- `deterministic-sim`: stable event ordering, seeded randomness, replay-safe transitions.
- `ostep-kernel`: OSTEP-aligned process, scheduling, syscall, memory, and interrupt modeling.
- `frontend-react-ts`: typed React UI for controls/visualizations with immutable DTO consumption.

## Suggested Usage Order

1. Start with `repo-workflow`.
2. Add stack/domain skill (`go-core`, `frontend-react-ts`, `ostep-kernel`, or `deterministic-sim`).
3. Apply `testing-discipline` before final verification.

## UI Design Baseline

For upcoming web work, frontend output should be simple, modern, and classic:

- restrained visual language with clear hierarchy and spacing,
- neutral palette with a single accent,
- typography chosen for readability over novelty,
- minimal, meaningful motion tied to simulator state changes.
