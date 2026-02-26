# Web UI

React + TypeScript frontend scaffold for realtime simulator control.

## Development

```bash
pnpm --dir=web install
pnpm --dir=web run dev
```

By default in dev, the app uses same-origin API calls and Vite proxies `/sessions`, `/lessons`, `/challenges`, and `/ws` to `http://127.0.0.1:8080`.

## UI Modes

- `/sandbox`: free-form simulator control and visualization
- `/challenge`: command-driven challenge attempts (`start -> interact -> check`) with deterministic grading and hint feedback

Challenge mode runs small OSTEP steps; Sandbox mode is for free experimentation.

## Verification

```bash
pnpm --dir=web run lint
pnpm --dir=web run typecheck
pnpm --dir=web run test
pnpm --dir=web run build
```
