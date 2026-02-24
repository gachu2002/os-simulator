# Milestone 05: Virtual Memory

## Goal

Implement address translation and page fault behavior.

## Scope

- VA to PA translation
- TLB model
- Page fault handler
- Replacement policy integration

## Deliverables

- Memory views for address space and frames
- Fault timeline visualization

## Exit Criteria

- Fault counts and frame occupancy match expected test cases

## Key Risks

- Adding too much MMU realism too early

## Suggested First Tasks

1. Implement baseline page table translation.
2. Add TLB with deterministic replacement.
3. Handle not-present and permission faults.
4. Add deterministic memory pressure tests.
