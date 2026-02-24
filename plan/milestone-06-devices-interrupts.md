# Milestone 06: Devices and Interrupts

## Goal

Model realistic asynchronous I/O and interrupt-driven wakeups.

## Scope

- Disk device
- Terminal device
- IRQ controller

## Deliverables

- I/O request and completion visualization
- Interrupt markers in timeline/trace

## Exit Criteria

- Blocked to ready wakeup traces are deterministic and correct

## Key Risks

- Race-like event ordering bugs when queue rules are weak

## Suggested First Tasks

1. Define device request lifecycle and completion interrupts.
2. Implement IRQ injection into event queue.
3. Connect blocked process wakeup path.
4. Add integration tests for syscall to interrupt to wakeup flow.
