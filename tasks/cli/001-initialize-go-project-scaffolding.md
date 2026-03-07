---
id: "001"
title: "Initialize Go project scaffolding"
status: completed
priority: high
dependencies: []
tags: ["go", "setup"]
created: 2026-03-07
---

# Initialize Go project scaffolding

## Objective

Set up the Go project structure for the viewmd Python-to-Go migration. This is the foundation all other tasks build on.

## Tasks

- [x] Run `go mod init github.com/driangle/viewmd/apps/cli`
- [x] Create directory structure: `apps/cli/cmd/viewmd/`, `apps/cli/internal/handler/`, `apps/cli/internal/frontmatter/`, `apps/cli/internal/render/`, `apps/cli/internal/classify/`
- [x] Create minimal `apps/cli/cmd/viewmd/main.go` entry point
- [x] Add placeholder `.go` files in each `apps/cli/internal/` package
- [x] Add `github.com/yuin/goldmark` dependency
- [x] Create `apps/cli/Makefile` with `build`, `test`, and `lint` targets
- [x] Create root `Makefile` delegating to `apps/cli/`

## Acceptance Criteria

- `go build ./...` succeeds
- `go test ./...` succeeds (with placeholder tests)
- `make build` produces a binary
- Directory structure matches the planned architecture
