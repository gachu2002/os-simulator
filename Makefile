SHELL := /bin/bash

DB_URL ?= postgres://postgres:postgres@localhost:5432/os_simulator?sslmode=disable
MIGRATE := go run github.com/golang-migrate/migrate/v4/cmd/migrate@latest
SQLC := go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest
AIR := go run github.com/air-verse/air@latest
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest
SHADCN := pnpm -C web dlx shadcn@latest

.PHONY: fmt lint test test-race test-deterministic test-coverage lesson-pack web-format-check web-lint web-typecheck web-test web-build security ci release-check sqlc-generate db-up db-down db-status db-create dev-server web-shadcn-add

fmt:
	gofmt -w .

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

ci: lint test test-deterministic test-race test-coverage lesson-pack web-lint web-typecheck web-test web-build

release-check: ci
	go build ./...

sqlc-generate:
	$(SQLC) generate -f sqlc.yaml

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

web-shadcn-add:
	@test -n "$(name)" || (printf "usage: make web-shadcn-add name=button\n" && exit 1)
	$(SHADCN) add --yes "$(name)"
