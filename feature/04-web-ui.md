# 04) Web UI

React + TypeScript challenge interface for OSTEP lesson-stage workflow.

```mermaid
flowchart TB
  APP[Challenge App] --> CHAL[/challenge]
  CHAL --> LESSON[Lesson Stage Runner]
  CHAL --> CVIZ[Challenge Snapshot]
  CHAL --> VIZ
  VIZ --> TL[Scheduler Timeline]
  VIZ --> MEM[Memory Panel]
  VIZ --> Q[Process Queues]
  VIZ --> PM[Process Metrics]
```

```mermaid
flowchart LR
  HTTP[createSession/fetchLessons/startChallenge/gradeChallenge] --> STATE[Reducer + Selectors + Query]
  WS[session events] --> STATE
  RUN[Challenge run output] --> MAP[challenge snapshot mapper]
  MAP --> STATE
  STATE --> UI[Deterministic Render]
```
