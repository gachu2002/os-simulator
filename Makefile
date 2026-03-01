SHELL := /bin/bash

DB_URL ?= postgres://postgres:postgres@localhost:5432/os_simulator?sslmode=disable
MIGRATE := go run github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1
SQLC := go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0
AIR := go run github.com/air-verse/air@v1.64.5
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

.PHONY: fmt fmt-check lint test test-race test-deterministic test-coverage lesson-pack web-format-check web-lint web-typecheck web-test web-build security audit-unused ci-go ci-web ci-security ci release-check sqlc-generate sqlc-verify db-up db-down db-status db-create dev-server

fmt:
	gofmt -w .

fmt-check:
	@files=$(gofmt -l .); \
	if [ -n "$$files" ]; then \
		echo "The following files are not gofmt-formatted:"; \
		echo "$$files"; \
		exit 1; \
	fi

lint:
	$(GOLANGCI_LINT) run

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

web-format-check:
	pnpm --dir=web exec prettier --check .

web-lint:
	pnpm --dir=web run lint

web-typecheck:
	pnpm --dir=web run typecheck

web-test:
	pnpm --dir=web run test

web-build:
	pnpm --dir=web run build

security:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	pnpm --dir=web audit --prod --audit-level high

audit-unused:
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...
	pnpm --dir=web exec tsc --noEmit --noUnusedLocals --noUnusedParameters

ci-go: fmt-check lint sqlc-verify test test-deterministic test-race test-coverage lesson-pack

ci-web: web-lint web-typecheck web-test web-build

ci-security: security

ci: ci-go ci-web

release-check: ci
	go build ./...

sqlc-generate:
	$(SQLC) generate -f sqlc.yaml

sqlc-verify: sqlc-generate
	git diff --exit-code -- internal/db/sqlc

db-up:
	$(MIGRATE) -path db/migrations -database "$(DB_URL)" up

db-down:
	$(MIGRATE) -path db/migrations -database "$(DB_URL)" down 1

db-status:
	$(MIGRATE) -path db/migrations -database "$(DB_URL)" version

db-create:
	@test -n "$(name)" || (printf "usage: make db-create name=create_table\n" && exit 1)
	$(MIGRATE) create -ext sql -dir db/migrations -seq "$(name)"

dev-server:
	$(AIR) -c .air.toml
