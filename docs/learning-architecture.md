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

The user experience is organized into three explicit modes:

- Path: guided curriculum and progression.
- Sandbox: free simulation and experimentation.
- Challenge: constrained assessments with grading.

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
- `objective`: what the learner should understand after completion.
- `prerequisites`: stage keys required before attempt.
- `difficulty`: intro, core, advanced.
- `estimated_minutes`: expected completion time.
- `concept_tags`: key concepts (for weak-spot analytics).
- `prompt`: the learner-facing task statement.
- `prediction_prompt`: required pre-run prediction prompt.
- `commands`: deterministic simulator command sequence.
- `validators`: objective checks.
- `explain_prompt`: required post-run reflection question.
- `hints`: nudge, concept, explicit escalation.
- `unlocks`: next stage keys enabled on pass.

This extends current `Lesson` and `Stage` structures and keeps deterministic command/validator behavior intact.

## Mastery Rules

- Pass threshold: mission validators all pass.
- Mastery threshold: pass plus low hint usage (L1-L2 preferred) and stable repeat success.
- Module mastery: all core missions passed, plus at least one challenge mission.
- Course completion: all modules mastered in sequence.

## Success Metrics

Track at mission and module levels:

- Attempt count before pass.
- Highest hint level used.
- Time-to-pass.
- Determinism stability check status.
- Concept weak spots (by `concept_tags`).

## Definition of Done for Product Spine

This step is complete when:

- All future learning features reference this document.
- New lessons are written against the mission contract.
- UI/transport work maps to one of the three product modes.
