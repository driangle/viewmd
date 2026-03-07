# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

viewmd is a tool for viewing Markdown files in the browser. It started as a Python CLI (`viewmd.py`) and is being rewritten in Go (`apps/cli/`). Both coexist during migration.

## Monorepo Structure

```
viewmd.py                  — Python CLI (legacy, still active)
apps/cli/                  — Go CLI (in development)
  cmd/viewmd/main.go       — Entry point
  internal/handler/        — HTTP handler and routing
  internal/frontmatter/    — Frontmatter parser
  internal/render/         — HTML template rendering
  internal/classify/       — File type classification
  Makefile                 — Go-specific build targets
tasks/cli/                 — Go migration task files
docs/plan/                 — Migration plan
```

## Commands

### Root (monorepo)

- **Build Go CLI:** `make build`
- **Run all Go tests:** `make test`
- **Lint everything:** `make lint` (runs `go vet` + `pylint`)
- **Clean build artifacts:** `make clean`
- **Set up git hooks:** `make setup`

### Python

- **Install for development:** `pip install -e ".[dev]"`
- **Run the server:** `viewmd` (default port 8000) or `viewmd <port>`
- **Run all tests:** `pytest`
- **Run a single test:** `pytest tests/test_frontmatter.py::TestParseFrontmatter::test_basic`
- **Build package:** `python -m build`
- **Lint:** `pylint viewmd.py tests/`

### Go CLI (`apps/cli/`)

- **Build:** `make -C apps/cli build` (or `make build` from root)
- **Test:** `make -C apps/cli test`
- **Lint:** `make -C apps/cli lint`

## Architecture

### Python (`viewmd.py`)

The entire application lives in `viewmd.py` — a single module with no internal package structure:

- `parse_frontmatter(content)` — Extracts YAML-like frontmatter (key: value pairs between `---` delimiters) from markdown text. Returns `(dict, body)` or `(None, original_content)`.
- `MarkdownHandler(SimpleHTTPRequestHandler)` — The HTTP request handler. Routes requests through `do_GET()` to one of three rendering paths:
  - `serve_markdown()` — Converts `.md`/`.markdown` files to HTML using the `markdown` library (with fenced_code, tables, nl2br extensions). Strips frontmatter and renders it as a styled HTML table above the content.
  - `serve_text_file()` — Wraps text files in a `<pre>` block with HTML escaping. Falls back to binary serving on `UnicodeDecodeError`.
  - `serve_directory_listing()` — Generates an HTML directory listing; auto-redirects to `README.md` if present.
- `main()` — Entry point. Binds `HTTPServer` on the given port with graceful error on port conflicts.

The version string is maintained in two places: `viewmd.py:VERSION` and `pyproject.toml:[project].version`.

### Go CLI (`apps/cli/`)

- `cmd/viewmd/main.go` — Entry point, CLI arg parsing, server startup
- `internal/handler/` — HTTP handler (routing, response helpers)
- `internal/frontmatter/` — Frontmatter parser (equivalent to Python `parse_frontmatter`)
- `internal/render/` — HTML template rendering (markdown, text, directory pages)
- `internal/classify/` — File type classification (text extensions, known filenames, dotfiles)

Uses `github.com/yuin/goldmark` for Markdown rendering and stdlib for HTTP/templates.

## Coding Guidelines

- **Max 200 lines per file.** If a file grows beyond this, split it into modules.
- **Max ~60 lines per function.** Extract helpers or break logic into smaller functions.
- **Add unit tests for non-trivial logic.** Use table-driven tests in Go. Focus on behavior and edge cases, not implementation details.
- Python lint enforced by pylint (`max-module-lines=200`, `max-statements=30`), configured in `pyproject.toml`.
- Go lint via `go vet`.

## Testing

### Python

Tests are in `tests/test_frontmatter.py` with fixtures in `tests/fixtures/`. Integration tests spin up a real HTTP server on a random port (`HTTPServer` bound to port 0) in a temp directory, copying fixture `.md` files there.

### Go

Tests live alongside source files (`*_test.go`). Use table-driven tests. Integration tests should use `net/http/httptest` or a real server on `:0` with temp directories for fixtures.
