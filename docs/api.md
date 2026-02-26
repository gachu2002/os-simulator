# API Contract

## Base Endpoints

- `GET /healthz`
- `POST /sessions`
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

## Session Create

### `POST /sessions`

Creates deterministic in-memory session.

Response: `201` with initial snapshot payload.

## Lessons List

### `GET /lessons`

Returns lesson summaries with lightweight stage metadata (`index`, `id`, `title`). Default catalog currently ships 20 lessons with 3 stages each.

## Challenge Start

### `POST /challenges/start`

Body:

```json
{
  "lesson_id": "l01-sched-rr-basics",
  "stage_index": 0
}
```

Returns challenge attempt metadata (`attempt_id`, `session_id`, objective, allowed commands, and limits).

## Challenge Grade

### `POST /challenges/grade`

Body:

```json
{
  "attempt_id": "a-000001"
}
```

Grades the current challenge session state and returns pass/fail, hint progression fields, output snapshot fields, and completion analytics.

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
