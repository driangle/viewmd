---
id: "007"
title: "Port integration tests to Go"
status: completed
priority: medium
dependencies: ["006"]
tags: ["go", "testing"]
created: 2026-03-07
---

# Port integration tests to Go

## Objective

Port all Python integration tests (test_frontmatter.py, test_server.py) to Go, using real HTTP servers on random ports with temp directories and fixture files.

## Tasks

- [ ] Set up test helper: start server on `:0`, copy fixtures to temp dir
- [ ] Port `TestServeMarkdownFrontmatter` tests (frontmatter table, no-frontmatter, HTML escaping, stripping)
- [ ] Port `TestServeTextFile` tests (pre wrapping, content escaping, dotfiles, known filenames, binary fallback)
- [ ] Port `TestServeDirectoryListing` tests (root listing, subdirs, README redirect, parent link)
- [ ] Port `TestRouting` tests (markdown routing, binary passthrough, 404, URL encoding)
- [ ] Copy fixture `.md` files into `testdata/` directory

## Acceptance Criteria

- All Python integration test scenarios covered in Go
- Tests use real HTTP servers on random ports (not mocks)
- Fixture files managed in `testdata/`
- `go test ./...` passes with all integration tests

## Sub-tasks

- **013** — Set up Go integration test infrastructure
- **014** — Port markdown and frontmatter integration tests
- **015** — Port text, directory, and routing integration tests
