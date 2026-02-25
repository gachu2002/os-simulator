# Web UI

React + TypeScript frontend scaffold for realtime simulator control.

## Development

```bash
pnpm --dir=web install
pnpm --dir=web run dev
```

By default, the app expects transport server at `http://localhost:8080`.

## UI Modes

- `/path`: guided lesson mode with lesson replay comparison
- `/sandbox`: free-form simulator control and visualization
- `/challenge`: live control + lesson goal validation
- `/progress`: persisted completion analytics and weak-concept summary

`/progress` now reads `GET /lessons/progress` for persisted analytics and weak-concept signals.

## Verification

```bash
pnpm --dir=web run lint
pnpm --dir=web run typecheck
pnpm --dir=web run test
pnpm --dir=web run build
```
