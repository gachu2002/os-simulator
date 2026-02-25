# Release Checklist (v1.0 Candidate)

## Pre-Release

- [ ] `make fmt`
- [ ] `make lint`
- [ ] `make test`
- [ ] `make test-race`
- [ ] `make test-coverage`
- [ ] `make test-deterministic`
- [ ] `make lesson-pack`
- [ ] `make security`
- [ ] `make release-check`
- [ ] CI workflow `ci` green on main branch head

## Artifact Validation

- [ ] `simcli` binary runs a basic syscall scenario
- [ ] replay hash remains stable for golden scenario
- [ ] lesson pack analytics reports 20/20 completion in smoke run
- [ ] `cmd/server` starts and responds on `/healthz`
- [ ] web build artifact (`web/dist`) is generated and deployable
- [ ] lesson runner API paths (`/lessons`, `/lessons/run`) respond successfully

## Deployment Smoke (Hosted)

- [ ] Backend hosted URL passes `/healthz`, `/lessons`, and `/lessons/run` checks
- [ ] Frontend hosted URL serves the app shell successfully
- [ ] GitHub workflow `deploy-smoke` passes on `main`

## Tagging and Notes

- [ ] Create annotated tag `v1.0.0-rc1`
- [ ] Add release notes: scope, known limitations, migration notes
- [ ] Include deterministic regression status and commit SHA in notes

## Post-Release

- [ ] Open next-iteration planning issue with lessons/content quality feedback
- [ ] Capture pilot checklist outcomes and completion analytics snapshot
