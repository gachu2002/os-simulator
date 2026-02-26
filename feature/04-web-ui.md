# 04) Web UI

React + TypeScript control plane organized into two learning modes: Sandbox and Challenge.

```mermaid
flowchart TB
  APP[App Router by URL mode] --> SBX[/sandbox]
  APP --> CHAL[/challenge]
  SBX --> CTRL[Control Bar]
  SBX --> STATUS[Status Cards]
  SBX --> LOG[Event Log]
  CHAL --> LESSON[Challenge Runner]
  CHAL --> CVIZ[Challenge Snapshot]
  SBX --> VIZ
  VIZ --> TL[Scheduler Timeline]
  VIZ --> MEM[Memory Panel]
  VIZ --> Q[Process Queues]
  VIZ --> PM[Process Metrics]
```

```mermaid
flowchart LR
  HTTP[createSession/fetchLessons/runLesson] --> STATE[Reducer + Selectors + Query]
  WS[session events] --> STATE
  RUN[Challenge run output] --> MAP[challenge snapshot mapper]
  MAP --> STATE
  STATE --> UI[Deterministic Render]
```
