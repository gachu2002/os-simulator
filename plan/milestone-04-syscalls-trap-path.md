# Milestone 04: Syscalls and Trap Path

## Goal

Make user-to-kernel control flow explicit and teachable.

## Scope

- Syscall dispatcher
- Trap entry and return path
- Minimal syscall set

## Deliverables

- `open`, `read`, `write`, `sleep`, `exit`
- Trace-friendly syscall path representation

## Exit Criteria

- Sequence trace assertions validate trap-save-dispatch-return ordering

## Key Risks

- Leaky abstraction between CPU and kernel layers

## Suggested First Tasks

1. Define syscall IDs and argument validation rules.
2. Implement trap frame save/restore flow.
3. Add minimal syscall handlers.
4. Add ordered trace assertions for every syscall.
