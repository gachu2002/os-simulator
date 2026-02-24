# Free Deployment + CD (Recommended)

This guide uses:

- Backend (`cmd/server`): Render free web service
- Frontend (`web`): Cloudflare Pages free tier
- CI/CD: GitHub Actions (`ci` + deploy smoke workflow)

## 1) Backend on Render (Go + WebSocket)

1. Create a new **Web Service** in Render from this GitHub repo.
2. Configure:
   - **Runtime**: Go
   - **Branch**: `main`
   - **Build Command**: `go build -o server ./cmd/server`
   - **Start Command**: `./server -addr :$PORT`
3. Leave service type as free tier.
4. Deploy and copy the backend URL (example: `https://os-sim-api.onrender.com`).

Notes:

- Free Render services can sleep when idle (cold starts).
- WebSocket is supported by Render for this use case.

## 2) Frontend on Cloudflare Pages

1. Create a new **Pages** project from this GitHub repo.
2. Configure build settings:
   - **Framework preset**: Vite
   - **Build command**: `pnpm --dir web build`
   - **Build output directory**: `web/dist`
   - **Root directory**: repository root
3. Deploy.

Optional env var (recommended):

- In Pages, set `VITE_API_BASE_URL` to your Render backend URL.

If you keep manual URL entry in the UI, no env var is required.

## 3) GitHub Environment Variables

Set these repository **Variables** (Settings -> Secrets and variables -> Actions -> Variables):

- `BACKEND_BASE_URL` = Render URL (no trailing slash)
- `FRONTEND_BASE_URL` = Pages URL (no trailing slash)

Example:

- `BACKEND_BASE_URL=https://os-sim-api.onrender.com`
- `FRONTEND_BASE_URL=https://os-sim-ui.pages.dev`

## 4) Deployment Smoke Workflow

Workflow file: `.github/workflows/deploy-smoke.yml`

It verifies:

- backend health endpoint: `GET /healthz`
- backend lesson list: `GET /lessons`
- backend lesson run: `POST /lessons/run`
- frontend availability: `GET /`

Run it manually with **Actions -> deploy-smoke -> Run workflow**, or let it run on pushes to `main`.

## 5) Suggested Free-Tier CD Policy

1. PR opens/updates -> run CI (`.github/workflows/ci.yml`).
2. Merge to `main` -> host auto-deploys from Git integration.
3. After merge -> run `deploy-smoke` workflow.
4. If smoke fails -> rollback by reverting merge commit.

## 6) Local Verification Before Merge

```bash
go test ./...
pnpm --dir web install
pnpm --dir web eslint .
pnpm --dir web tsc --noEmit
pnpm --dir web vitest run
pnpm --dir web build
```
