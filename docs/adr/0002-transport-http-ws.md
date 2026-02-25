# ADR 0002: HTTP + WebSocket Transport

## Status

Accepted

## Decision

Use HTTP for session/lesson request-response operations and WebSocket for realtime command/event streaming.

## Rationale

This split keeps control plane simple while preserving low-latency event flow for interactive UI.

## Consequences

- API contracts require explicit versioning and typed DTOs.
- CORS/origin controls and WS origin checks are mandatory for hosted deployments.
