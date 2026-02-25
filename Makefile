SHELL := /bin/bash

.PHONY: fmt lint test test-race test-deterministic test-coverage lesson-pack web-lint web-typecheck web-test web-build security ci release-check

fmt:
	gofmt -w .

lint:
	golangci-lint run

test:
	go test ./...

test-race:
	go test -race ./...

test-coverage:
	bash scripts/ci/check_go_coverage.sh

test-deterministic:
	go test ./internal/sim -run 'TestGoldenTraceHash|TestReplayFromLogMatchesOriginalHash|TestSyscallToIRQToWakeupFlowIsDeterministic'
	go test ./internal/lessons -run 'TestScenarioLessonsPassWithExpectedFeedbackKeys|TestCompletionAnalyticsAndPilotChecklist'

lesson-pack:
	go run ./cmd/simcli -run-lesson-pack

web-lint:
	pnpm --dir web eslint .

web-typecheck:
	pnpm --dir web tsc --noEmit

web-test:
	pnpm --dir web vitest run

web-build:
	pnpm --dir web build

security:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	pnpm --dir web audit --prod --audit-level high

ci: lint test test-deterministic test-race test-coverage lesson-pack web-lint web-typecheck web-test web-build

release-check: ci
	go build ./...
