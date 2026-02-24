SHELL := /bin/bash

.PHONY: fmt test test-deterministic lesson-pack web-lint web-typecheck web-test web-build ci release-check

fmt:
	gofmt -w .

test:
	go test ./...

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

ci: test test-deterministic lesson-pack web-lint web-typecheck web-test web-build

release-check: ci
	go build ./...
