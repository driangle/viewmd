---
id: "001"
title: "Initialize Go project scaffolding"
status: pending
priority: high
dependencies: []
tags: ["go", "setup"]
created: 2026-03-07
---

# Initialize Go project scaffolding

## Objective

Set up the Go project structure for the viewmd Python-to-Go migration. This is the foundation all other tasks build on.

## Tasks

- [ ] Run `go mod init github.com/driangle/viewmd`
- [ ] Create directory structure: `cmd/viewmd/`, `internal/handler/`, `internal/frontmatter/`, `internal/render/`, `internal/classify/`
- [ ] Create minimal `cmd/viewmd/main.go` entry point
- [ ] Add placeholder `.go` files in each `internal/` package
- [ ] Add `github.com/yuin/goldmark` dependency
- [ ] Create `Makefile` with `build`, `test`, and `lint` targets

## Acceptance Criteria

- `go build ./...` succeeds
- `go test ./...` succeeds (with placeholder tests)
- `make build` produces a binary
- Directory structure matches the planned architecture
