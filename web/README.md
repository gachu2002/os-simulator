# Web UI (Milestone 12)

React + TypeScript frontend scaffold for realtime simulator control.

## Development

```bash
pnpm --dir web install
pnpm --dir web dev
```

By default, the app expects transport server at `http://localhost:8080`.

## Verification

```bash
pnpm --dir web eslint .
pnpm --dir web tsc --noEmit
pnpm --dir web vitest run
pnpm --dir web build
```
