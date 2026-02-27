# OSTEP Chapter Map v2

This file is the canonical chapter coverage map for curriculum iteration.

Status meanings:

- `engine-ready`: simulator supports challenge-level evaluation for this chapter focus.
- `engine-partial`: simulator can teach core ideas but not full chapter breadth.
- `theory-only`: chapter content can be taught, but challenge support requires new simulator capability.

## Capability Baseline (Current Simulator)

Supported core mechanics now:

- CPU schedulers: `fifo`, `rr`, `mlfq`
- process trace events: dispatch, compute, preempt, wakeup, exit, trap enter/return
- syscall path: `open`, `read`, `write`, `sleep`, `exit`
- memory model: address translation, paging faults, FIFO replacement, TLB counters
- filesystem model: path traversal, block mapping, invariants check
- deterministic challenge grading: trace, metrics, faults, fs invariants

Missing mechanics for full OSTEP parity:

- lottery scheduling and multi-CPU behavior model
- segmentation and allocator policy depth
- advanced page tables and complete VM system behavior
- lock/CV/semaphore/event-loop runtime primitives
- crash/recovery transaction modeling depth (journal replay modes)
- distributed FS consistency protocols (NFS/AFS)
- security model primitives (authn/authz/crypto/distributed security)

## Chapter-to-Curriculum Coverage Map

## Intro

- `2 Introduction` -> `engine-ready` (theory + trace-level challenge feasible)
- `4 Processes` -> `engine-ready` (covered by current process/scheduler flow)
- `5 Process API` -> `engine-partial` (syscalls are simplified)
- `6 Direct Execution` -> `engine-ready` (trap path visible in trace)

## CPU Virtualization

- `7 CPU Scheduling` -> `engine-ready`
- `8 Multi-level Feedback` -> `engine-ready`
- `9 Lottery Scheduling` -> `theory-only` (policy implementation missing)
- `10 Multi-CPU Scheduling` -> `theory-only` (multi-CPU model missing)

## Memory Virtualization

- `13 Address Spaces` -> `engine-ready`
- `14 Memory API` -> `engine-partial` (user API details simplified)
- `15 Address Translation` -> `engine-ready`
- `16 Segmentation` -> `theory-only` (segmentation model missing)
- `17 Free Space Management` -> `theory-only` (allocator strategies missing)
- `18 Introduction to Paging` -> `engine-ready`
- `19 TLBs` -> `engine-ready`
- `20 Advanced Page Tables` -> `theory-only`
- `21 Swapping: Mechanisms` -> `theory-only`
- `22 Swapping: Policies` -> `theory-only`
- `23 Complete VM Systems` -> `theory-only`

## Concurrency

- `26 Concurrency and Threads` -> `engine-partial` (single-process blocking visible; threads not modeled)
- `27 Thread API` -> `theory-only`
- `28 Locks` -> `theory-only`
- `29 Locked Data Structures` -> `theory-only`
- `30 Condition Variables` -> `theory-only`
- `31 Semaphores` -> `theory-only`
- `32 Concurrency Bugs` -> `theory-only`
- `33 Event-based Concurrency` -> `theory-only`

## Persistence

- `36 I/O Devices` -> `engine-ready`
- `37 Hard Disk Drives` -> `engine-partial` (latency modeled; full disk model simplified)
- `38 RAID` -> `theory-only`
- `39 Files and Directories` -> `engine-ready`
- `40 File System Implementation` -> `engine-ready`
- `41 FFS` -> `theory-only`
- `42 FSCK and Journaling` -> `engine-partial` (invariants yes, full journal replay modes missing)
- `43 LFS` -> `theory-only`
- `44 SSDs` -> `theory-only`
- `45 Data Integrity and Protection` -> `engine-partial`

## Distributed Systems

- `48 Distributed Systems` -> `theory-only`
- `49 NFS` -> `theory-only`
- `50 AFS` -> `theory-only`

## Security

- `53 Intro Security` -> `theory-only`
- `54 Authentication` -> `theory-only`
- `55 Access Control` -> `theory-only`
- `56 Cryptography` -> `theory-only`
- `57 Distributed Security` -> `theory-only`

## Current Lesson Mapping (Active IDs)

- CPU scheduling lessons: `l01` to `l06c` -> chapters 7-8 coverage (`engine-ready`)
- Memory lessons: `l07` to `l11c` -> chapters 13, 15, 18, 19 coverage (`engine-ready`)
- Concurrency-related blocking/IRQ lessons: `l12` to `l15c` -> chapter 26 + I/O bridge (`engine-partial`)
- Persistence lessons: `l16` to `l20c` -> chapters 36, 39, 40, partial 42/45 (`engine-ready` to `engine-partial`)

Not yet represented as active lessons:

- Intro chapter group (`2`, `5`, `6`) as first-class section
- chapters `9`, `10`, `16`, `17`, `20-23`, `27-33`, `38`, `41`, `43`, `44`, `48-50`, `53-57`

## Revision Notes (Step 1)

- The previous plan overestimated practical concurrency coverage; most lock/CV/semaphore content is currently theory-only.
- Lottery and multi-CPU must be flagged as explicit capability gaps, not silently implied by lesson titles.
- Journaling coverage must be marked partial until replay/recovery mode behavior exists in simulator core.

This map is the gate for Step 2 schema design and Step 3 lesson rewrites.
