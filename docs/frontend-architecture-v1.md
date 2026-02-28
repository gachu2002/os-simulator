# Frontend Architecture v1

## Goals

- Keep routing and navigation declarative with React Router.
- Keep server-state fetching in TanStack Query.
- Keep client/session workflow state in Zustand.
- Keep domain boundaries explicit: API transport -> model mapping -> feature UI.
- Keep feature modules isolated and easier to maintain.

## Layering

- `src/app`: app shell, providers, router, layout.
- `src/features`: user-facing feature slices (curriculum, learn, challenge, visualization).
- `src/entities`: shared domain model types (lesson, challenge, snapshot).
- `src/shared`: cross-cutting code (api client, config, utils, styles).

## State Ownership

- TanStack Query: curriculum and lesson-read API data.
- Zustand `sessionStore`: websocket session lifecycle and snapshot stream.
- Zustand `challengeStore`: attempt state, submission state, challenge UI controls.
- Local component state: only highly local, transient UI concerns.

## Routing Contract

- `/` -> curriculum overview.
- `/lesson/:lessonID/learn` -> theory page.
- `/lesson/:lessonID/challenge` -> interactive challenge page.

## Styling Contract

- Tailwind-first for layout and component composition.
- Legacy CSS retained only where visualization components still depend on existing class selectors.
- New UI code should prefer utility classes and shared `cn()` helper.

## Near-term Refactor Queue

1. Move legacy `components/*` pages into feature-local `ui/*` modules.
2. Add DTO-to-domain mappers to isolate backend payload changes.
3. Split large challenge page into panel-specific components and hooks.
4. Gradually retire legacy CSS selectors in favor of Tailwind utility components.
