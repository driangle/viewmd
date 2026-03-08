# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

viewmd is a minimal file viewer for the browser, written in Go.

**Design principle: View everything, render what you can, manage nothing.**

viewmd is not a file manager — it doesn't upload, edit, rename, or delete. It's a read-only viewer that renders each file type as nicely as possible:
- **Markdown** — full rendered HTML (core value)
- **Code/text** — syntax-highlighted display
- **Images** — inline rendering
- **Directories** — clean listing with README auto-display
- **Other files** — listed in directories, download or "can't preview" on open

## Project Structure

```
apps/cli/                  — Go CLI
  cmd/viewmd/main.go       — Entry point
  internal/handler/        — HTTP handler and routing
  internal/frontmatter/    — Frontmatter parser
  internal/render/         — HTML template rendering
  internal/classify/       — File type classification
  Makefile                 — Go build targets
tasks/cli/                 — Task files
python/                    — Archived Python version (PyPI-published)
```

## Commands

- **Build:** `make build`
- **Test:** `make test`
- **Lint:** `make lint`
- **Clean:** `make clean`
- **Install:** `make install`
- **Set up git hooks:** `make setup`

## Architecture

- `cmd/viewmd/main.go` — Entry point, CLI arg parsing, server startup
- `internal/handler/` — HTTP handler (routing, response helpers)
- `internal/frontmatter/` — Frontmatter parser
- `internal/render/` — HTML template rendering (markdown, text, directory pages)
- `internal/classify/` — File type classification (text extensions, known filenames, dotfiles)

Uses `github.com/yuin/goldmark` for Markdown rendering and stdlib for HTTP/templates.

## Coding Guidelines

- **Max 200 lines per file.** If a file grows beyond this, split it into modules.
- **Max ~60 lines per function.** Extract helpers or break logic into smaller functions.
- **Add unit tests for non-trivial logic.** Use table-driven tests. Focus on behavior and edge cases, not implementation details.
- Go lint via `go vet`.

## Documentation

When adding a new feature or making a significant change to existing behavior, update the relevant docs (`README.md`, `docs/reference/cli.md`, `docs/getting-started/quick-start.md`, etc.). New CLI flags, changed defaults, and user-visible features should always be documented.

## Testing

Tests live alongside source files (`*_test.go`). Use table-driven tests. Integration tests should use `net/http/httptest` or a real server on `:0` with temp directories for fixtures.
