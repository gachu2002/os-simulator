# Frontend Coding Standards v1

## TypeScript

- Strict mode stays enabled.
- Avoid `any`; prefer explicit interfaces at boundaries.
- Keep transport DTOs separate from domain-facing models.

## React Components

- Container components orchestrate data and state.
- Presentational components render UI only.
- Avoid combining fetch, websocket, and rendering logic in one component.

## State

- Query state belongs in TanStack Query hooks.
- Session/challenge workflow state belongs in Zustand stores.
- Local state is for local form/input interaction only.

## Styling

- Use Tailwind utility classes by default.
- Use `cn()` helper for conditional class composition.
- Avoid adding broad global selectors unless necessary for legacy compatibility.

## File and Import Conventions

- Prefer feature-local imports over cross-feature deep imports.
- Keep shared code under `src/shared` and domain model types under `src/entities`.
- Keep one primary responsibility per file.

## Testing

- Unit test selectors and store actions.
- Integration test challenge flow and route rendering.
- Keep test names behavior-oriented and deterministic.
