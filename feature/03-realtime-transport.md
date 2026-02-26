# 03) Realtime Transport

HTTP + WebSocket transport over deterministic simulator sessions.

```mermaid
flowchart LR
  FE[Frontend] -->|GET /lessons| API[cmd/server]
  FE -->|POST /challenges/start| API
  FE -->|POST /challenges/grade| API
  FE <-->|WS /ws/{sessionID}| API
  API --> SM[SessionManager]
  SM --> ENG[sim.Engine]
  API --> LE[lessons.Engine]
```

```mermaid
sequenceDiagram
  participant C as Client
  participant S as Server
  participant E as Session
  C->>S: POST /challenges/start
  S->>E: create deterministic engine
  S-->>C: attempt_id + session_id + limits
  C->>S: WS command envelope
  S->>E: Apply(command)
  S-->>C: session.snapshot or session.error
```
