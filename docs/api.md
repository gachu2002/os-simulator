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

Returns lesson summaries with stage metadata (`index`, `id`, `title`, `theory`, `objective`, `pass_conditions`, `prerequisites`, `allowed_commands`, `limits`) plus progress status (`attempts`, `completed`, `unlocked`). Use optional query param `learner_id` to scope unlock/progress per learner. Default catalog currently ships 20 lessons with 3 stages each.

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

Returns lesson-stage attempt metadata (`attempt_id`, `session_id`, objective, allowed commands, and limits).

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

For challenge sessions, `session.snapshot` may include `snapshot.challenge` budget fields (`max_steps`, `used_steps`, `remaining_steps`, and policy-change counterparts) to drive live exercise limits in the UI.
