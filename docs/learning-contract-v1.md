# Learning Product Contract v1

This document is the implementation contract for Phase 0.

It defines the non-negotiable UX, content model, and grading behavior for the OSTEP learning product.

## Product Promise

- The simulator teaches OSTEP through lesson-first learning, not free-form controls.
- Every lesson has exactly two pages: Learn and Challenge.
- Learn page contains theory only.
- Challenge page contains exactly three sections: Actions, Visualization, Goal+Submit.
- Challenge grading is deterministic and replay-safe.

## Route Contract

- `/`: course home with section cards and lesson titles.
- `/lesson/:lessonId/learn`: theory page for one lesson.
- `/lesson/:lessonId/challenge`: challenge page for one lesson.

No mixed learn/exercise mode within one page.

## Home Page Contract

- Home shows top-level OSTEP sections in sequence.
- Expanding a section reveals lesson titles only.
- Each lesson row shows:
  - title
  - status: locked, ready, passed
  - estimated minutes
- Clicking a lesson opens `/lesson/:lessonId/learn`.

## Learn Page Contract

Learn page has no simulator controls and no grading controls.

Each lesson Learn page must contain these blocks in order:

1. Core Idea
2. Mechanism
3. Worked Example
4. Common Mistakes
5. What To Watch In Challenge

Minimum quality bar:

- A learner can predict expected challenge behavior before running challenge actions.
- The theory covers all concepts required by challenge pass conditions.

## Challenge Page Contract

The page layout is fixed to exactly three sections.

### 1) Actions

- Show only allowed actions for this lesson.
- Hide disallowed actions.
- Show limits and remaining budgets:
  - steps
  - policy changes
  - config changes

### 2) Visualization

- Show only concept-relevant visual panels.
- Panels are deterministic views of simulator state.
- Required panel families:
  - timeline/trace
  - process/queue state
  - metrics
  - memory state (when memory lesson)
  - filesystem state (when persistence lesson)

### 3) Goal + Submit

- Human-readable goal statement.
- Explicit pass checklist.
- Submit action evaluates deterministic validators.
- Result block always includes:
  - Passed/Failed
  - expected vs actual for each check
  - first failed check highlighted
  - actionable hint mapped to failed check

## Grading Contract

Pass formula:

`PASS = all required checks passed AND no forbidden checks violated AND limits respected`

Failure formula:

`FAIL = any required check failed OR any forbidden check violated OR any limit exceeded`

### Validator Set

Supported now:

- `trace_contains_all`
- `trace_order`
- `trace_count_eq`
- `trace_count_lte`
- `no_event`
- `metric_eq`
- `metric_gte`
- `metric_lte`
- `fault_eq`
- `fault_lte`
- `fs_ok`

Planned extensions (Phase 2+):

- `budget_ok`

## Determinism Contract

- The lesson challenge seed is fixed per lesson.
- Replay with same action sequence must produce identical grade results.
- Equal-priority event ordering must remain stable.
- No wall-clock dependence in grading path.

## Content Model Contract

The curriculum is content-first. Runtime code consumes content definitions.

Current runtime source of truth:

- `internal/lessons/catalog_content_v1.json`
- `internal/lessons/lesson_content_v2/*.json` (theory and chapter-grounded learn content, one file per lesson)
- `internal/lessons/lesson_stage_content_v2/*.json` (objective/goal/hints, one file per lesson)

Canonical entities:

- Section
- Lesson
- LearnContent
- ChallengeSpec
- ValidatorSpec
- HintMap

Curriculum v2 authoring references:

- `docs/curriculum/ostep-chapter-map-v2.md`
- `docs/curriculum/lesson-schema-v2.md`
- `docs/curriculum/lesson-v2.schema.json`

### Section

- `id`
- `title`
- `subtitle`
- `order`
- `coming_soon`

### Lesson

- `id`
- `section_id`
- `title`
- `estimated_minutes`
- `difficulty`
- `chapter_refs`
- `prerequisites`
- `learn`
- `challenge`

### LearnContent

- `core_idea`
- `mechanism_steps[]`
- `worked_example`
- `common_mistakes[]`
- `challenge_watch[]`

### ChallengeSpec

- `objective`
- `goal`
- `bootstrap_commands[]`
- `allowed_actions[]`
- `limits`
- `pass_conditions[]`
- `expected_visual_cues[]`
- `validators[]`
- `hint_map[]`

## API Contract (Phase 0 target shape)

- `GET /curriculum`: section + lesson summaries for home.
- `GET /lessons/:lessonId/learn`: learn content only.
- `POST /challenges/start`: starts deterministic attempt.
- `POST /challenges/submit`: grades attempt and returns pass/fail details.

Migration note: existing endpoints may remain temporarily, but these are canonical.

## Completion Criteria for Phase 0

Phase 0 is complete when:

- This contract is accepted and referenced by implementation tasks.
- Every new lesson/challenge content entry conforms to this contract.
- UI implementation work can be validated directly against this document.
