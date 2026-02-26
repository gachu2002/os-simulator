# Learning Architecture (Step 1: Product Spine)

This document defines the product learning spine so simulator features map to a clear teaching journey.

## Learner Persona

- Primary learner: self-taught or early-career engineer learning OSTEP by doing.
- Starting point: knows basic programming, weak OS intuition.
- Goal: build reliable mental models for scheduling, memory, concurrency, and persistence.

## Core Learning Loop

Every mission must follow the same loop:

1. Learn concept
2. Predict behavior
3. Run simulation
4. Explain observed result
5. Check understanding
6. Unlock next mission

## Terminology

- Lesson: a curriculum unit within one module.
- Stage mission: one executable step inside a lesson.
- Stage key: `<lesson_id>:<stage_id>` prerequisite/unlock identifier.

## Product Modes

The user experience is organized into two explicit modes:

- Sandbox: free simulation and experimentation.
- Challenge: constrained challenge-step assessments with grading.

Each screen should belong to one mode only.

## Curriculum Order

Module unlock sequence:

1. CPU Virtualization and Scheduling
2. Memory Virtualization
3. Concurrency and Interrupts
4. Persistence and Filesystem

Prerequisites flow in this same order; later modules require baseline mastery in earlier modules.

## Mission Contract (Canonical Lesson Schema v2)

Every stage mission should define:

- `id`: stable stage id.
- `module`: curriculum module id.
- `title`: learner-facing step name.
- `prerequisites`: stage keys required before attempt.
- `commands`: deterministic simulator command sequence.
- `validators`: objective checks.
- `hints`: nudge, concept, explicit escalation.

This extends current `Lesson` and `Stage` structures and keeps deterministic command/validator behavior intact.

## Mastery Rules

- Pass threshold: mission validators all pass.
- Mastery threshold: pass all required stage missions for a module.
- Module mastery: all core missions passed, plus at least one challenge mission.
- Course completion: all modules mastered in sequence.

## Success Metrics

Track at mission and module levels:

- Attempt count before pass.
- Stage completion rate.
- Module completion rate.

## Definition of Done for Product Spine

This step is complete when:

- All future learning features reference this document.
- New lessons are written against the mission contract.
- UI/transport work maps to one of the two product modes.
