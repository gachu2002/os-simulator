# API Contract

## Base Endpoints

- `GET /healthz`
- `GET /lessons`
- `POST /challenges/start`
- `POST /challenges/grade`
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

## Lessons List

### `GET /lessons`

Returns lesson summaries with section metadata (`section_id`, `section_title`, `difficulty`, `estimated_minutes`, `chapter_refs`) and stage metadata (`index`, `id`, `title`, `theory`, `theory_detail`, `objective`, `goal`, `pass_conditions`, `prerequisites`, `allowed_commands`, `action_descriptions`, `expected_visual_cues`, `limits`) plus progress status (`attempts`, `completed`, `unlocked`). `allowed_commands` may include tuning actions such as `set_frames`, `set_tlb_entries`, `set_disk_latency`, and `set_terminal_latency`. `limits` includes `max_steps`, `max_policy_changes`, and `max_config_changes`. Use optional query param `learner_id` to scope unlock/progress per learner. Default catalog currently ships 28 lessons with 3 stages each.

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

## Challenge Grade

### `POST /challenges/grade`

Body:

```json
{
  "attempt_id": "a-000001",
  "learner_id": "learner-123"
}
```

Grades the current lesson-stage session state and returns pass/fail, hint progression fields, output snapshot fields, completion analytics, and per-validator check results in `validator_results`. `learner_id` must match the learner that started the attempt.

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
