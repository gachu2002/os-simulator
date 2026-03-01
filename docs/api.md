# API Contract (V3)

This document describes the active V3 API used by the frontend.

## Base Endpoints

- `GET /healthz`
- `GET /curriculum/v3`
- `GET /lessons/{lessonID}/learn/v3`
- `GET /lessons/{lessonID}/challenge/v3`
- `POST /challenges/start/v3`
- `POST /challenges/action/v3`
- `POST /challenges/submit/v3`
- `GET /challenges/{attemptID}/replay/v3`

## Request/Response Notes

- Content type: `application/json`
- Request body size capped to 1 MiB
- Unknown JSON fields are rejected on write endpoints

## Error Envelope

All HTTP errors return:

```json
{
  "code": "string",
  "message": "string",
  "request_id": "req-00000001"
}
```

`X-Request-ID` is echoed in response headers; clients may send their own.

## Curriculum

### `GET /curriculum/v3`

Returns Section 1 (`virtualization-cpu`) curriculum model with lesson theory concepts,
challenge actions, visualizer specs, optional challenge parts (`A`/`B`), and
cross-cutting feature flags.

## Lesson Learn

### `GET /lessons/{lessonID}/learn/v3`

Returns V3 lesson metadata for Section 1 (`version`, `section_id`, and lesson payload).

### `GET /lessons/{lessonID}/challenge/v3`

Returns V3 challenge manifest for a lesson, including:

- `actions`
- `action_capabilities.supported_now`
- `action_capabilities.planned`
- `action_capability_notes` (per-action status/reason/fallback guidance)
- `part_required` and optional `parts`
- `visualizer`
- `cross_cutting_features`

## Challenge Start

### `POST /challenges/start/v3`

Body:

```json
{
  "lesson_id": "l01-process-basics",
  "part_id": "A",
  "learner_id": "learner-123"
}
```

Starts a V3 challenge attempt and returns lesson metadata, optional part metadata,
runtime limits, action/visualizer descriptors, `action_capabilities`, and
`action_capability_notes`.

`part_id` is required only for lessons that define challenge parts.

## Challenge Action

### `POST /challenges/action/v3`

Applies one V3 action against an attempt and returns the mapped simulator event.

```json
{
  "attempt_id": "a-000001",
  "learner_id": "learner-123",
  "action": "execute_instruction",
  "count": 1
}
```

Response includes `mapped_command` and `event` (`session.snapshot` or `session.error`).
Unsupported actions return `400` with `code="invalid_action"`.

## Challenge Submit

### `POST /challenges/submit/v3`

Body:

```json
{
  "attempt_id": "a-000001",
  "learner_id": "learner-123"
}
```

Grades a V3 attempt and returns V3 lesson/part metadata, pass/fail + hint fields,
output snapshot, analytics, validator results, and `action_capabilities`.

Top-level response contract fixtures used by tests:

- `contracts/challenge_start_v3_contract.json`
- `contracts/challenge_submit_v3_contract.json`

## Challenge Replay

### `GET /challenges/{attemptID}/replay/v3`

Returns replay payload for a V3 attempt including full trace, trace hash/length,
process/metrics/memory snapshots, filesystem status, and `action_capabilities`.

`learner_id` query parameter is required and must match the attempt owner.
