# Web UI

React + TypeScript frontend scaffold for realtime simulator control.

## Development

```bash
pnpm --dir=web install
pnpm --dir=web run dev
```

By default in dev, the app uses same-origin API calls and Vite proxies `/healthz`, `/curriculum`, `/lessons`, and `/challenges` to `http://127.0.0.1:8080`.

## UI Modes

- `/`: section-first course overview
- `/lesson/:lessonId/learn`: theory-only lesson page
- `/lesson/:lessonId/challenge`: challenge page (`actions -> visualization -> goal + submit`)

Challenge mode runs focused OSTEP lesson challenges.

## Verification

```bash
pnpm --dir=web run lint
pnpm --dir=web run typecheck
pnpm --dir=web run test
pnpm --dir=web run build
```
