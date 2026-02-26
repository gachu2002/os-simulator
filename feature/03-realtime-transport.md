# 03) Realtime Transport

HTTP + WebSocket transport over deterministic simulator sessions.

```mermaid
flowchart LR
  FE[Frontend] -->|POST /sessions| API[cmd/server]
  FE -->|GET /lessons| API
  FE -->|POST /lessons/run| API
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
  C->>S: POST /sessions
  S->>E: create deterministic engine
  S-->>C: session_id + initial snapshot
  C->>S: WS command envelope
  S->>E: Apply(command)
  S-->>C: session.snapshot or session.error
```
