SHELL := /bin/bash

.PHONY: fmt test test-deterministic lesson-pack ci release-check

fmt:
	gofmt -w .

test:
	go test ./...

test-deterministic:
	go test ./internal/sim -run 'TestGoldenTraceHash|TestReplayFromLogMatchesOriginalHash|TestSyscallToIRQToWakeupFlowIsDeterministic'
	go test ./internal/lessons -run 'TestScenarioLessonsPassWithExpectedFeedbackKeys|TestCompletionAnalyticsAndPilotChecklist'

lesson-pack:
	go run ./cmd/simcli -run-lesson-pack

ci: test test-deterministic lesson-pack

release-check: ci
	go build ./...
