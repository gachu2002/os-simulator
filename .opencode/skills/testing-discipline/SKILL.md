---
name: testing-discipline
description: Add focused deterministic tests, run targeted test filters first, and expand to full suites after local behavior is proven.
compatibility: opencode
metadata:
  focus: quality
  priority: high
---

## What I do
- Add regression tests alongside behavior changes.
- Prefer fast targeted tests before full suite execution.
- Keep tests deterministic and reproducible.

## When to use me
- Use for every non-trivial code change.
- Use immediately when fixing a bug or flaky behavior.

## Single-test patterns
- Go single test: `go test ./path/to/pkg -run '^TestName$'`
- Go subtest: `go test ./path/to/pkg -run 'TestName/sub_case'`
- Vitest by file: `pnpm --dir web vitest run src/path/file.test.ts`
- Vitest by name: `pnpm --dir web vitest run -t "case name"`
- Playwright by title: `pnpm --dir web playwright test -g "test title"`

## Verification sequence
- Run targeted tests first.
- Run affected package or test file suite.
- Finish with full suite for touched stack.
