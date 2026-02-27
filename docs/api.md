# API Contract

## Base Endpoints

- `GET /healthz`
- `GET /curriculum`
- `GET /lessons/{lessonID}/learn`
- `POST /challenges/start`
- `POST /challenges/submit`
- `WS /ws/{sessionID}`

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

### `GET /curriculum`

Returns section-first curriculum payload for home page rendering. Response includes ordered sections (`id`, `title`, `subtitle`, `order`, `coming_soon`) and embedded lesson summaries. For active sections, progress totals are included (`completed_stages`, `total_stages`). Supports optional `learner_id` query param.

## Lesson Learn

### `GET /lessons/{lessonID}/learn`

Returns theory-first lesson payload intended for Learn page rendering. Includes lesson metadata and per-stage learn blocks (`theory`, `theory_detail`, `objective`, `goal`, prerequisites, expected visual cues). This endpoint excludes challenge control metadata.

## Challenge Start

### `POST /challenges/start`

Body:

```json
{
  "lesson_id": "l01-sched-rr-basics",
  "stage_index": 0,
  "learner_id": "learner-123"
}
```

Returns lesson-stage attempt metadata (`attempt_id`, `session_id`, objective, allowed commands, and limits including optional `max_config_changes`).

Also includes `goal` and `pass_conditions` so the challenge page can render submit checklist before grading.

## Challenge Submit

### `POST /challenges/submit`

Body:

```json
{
  "attempt_id": "a-000001",
  "learner_id": "learner-123"
}
```

Grades the current lesson-stage session state and returns pass/fail, hint progression fields, output snapshot fields, completion analytics, and per-validator check results in `validator_results`. `learner_id` must match the learner that started the attempt.

Each validator result includes `expected` and `actual` fields for expected-vs-actual rendering.

## WebSocket Stream

### `WS /ws/{sessionID}`

Inbound envelope:

```json
{
  "type": "command",
  "command": { "name": "step", "count": 1 }
}
```

Outbound event types:

- `session.snapshot`
- `session.error`

For challenge sessions, `session.snapshot` may include `snapshot.challenge` budget fields (`max_steps`, `used_steps`, `remaining_steps`, plus policy and config-change counterparts) to drive live exercise limits in the UI.
