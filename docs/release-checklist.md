# Release Checklist (v1.0 Candidate)

## Pre-Release

- [ ] `make fmt`
- [ ] `make test`
- [ ] `make test-deterministic`
- [ ] `make lesson-pack`
- [ ] `make release-check`
- [ ] CI workflow `ci` green on main branch head

## Artifact Validation

- [ ] `simcli` binary runs a basic syscall scenario
- [ ] replay hash remains stable for golden scenario
- [ ] lesson pack analytics reports 20/20 completion in smoke run

## Tagging and Notes

- [ ] Create annotated tag `v1.0.0-rc1`
- [ ] Add release notes: scope, known limitations, migration notes
- [ ] Include deterministic regression status and commit SHA in notes

## Post-Release

- [ ] Open next-iteration planning issue with lessons/content quality feedback
- [ ] Capture pilot checklist outcomes and completion analytics snapshot
