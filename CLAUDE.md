# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

viewmd is a single-file Python CLI tool that serves as an HTTP server for viewing Markdown files in the browser. Installed via `pipx install viewmd`, it renders `.md` files as styled HTML, displays text files with syntax-appropriate formatting, and provides directory browsing. Published to PyPI via GitHub Actions on release.

## Commands

- **Install for development:** `pip install -e ".[dev]"`
- **Run the server:** `viewmd` (default port 8000) or `viewmd <port>`
- **Run all tests:** `pytest`
- **Run a single test:** `pytest tests/test_frontmatter.py::TestParseFrontmatter::test_basic`
- **Build package:** `python -m build`
- **Lint:** `pylint viewmd.py tests/`

## Architecture

The entire application lives in `viewmd.py` — a single module with no internal package structure:

- `parse_frontmatter(content)` — Extracts YAML-like frontmatter (key: value pairs between `---` delimiters) from markdown text. Returns `(dict, body)` or `(None, original_content)`.
- `MarkdownHandler(SimpleHTTPRequestHandler)` — The HTTP request handler. Routes requests through `do_GET()` to one of three rendering paths:
  - `serve_markdown()` — Converts `.md`/`.markdown` files to HTML using the `markdown` library (with fenced_code, tables, nl2br extensions). Strips frontmatter and renders it as a styled HTML table above the content.
  - `serve_text_file()` — Wraps text files in a `<pre>` block with HTML escaping. Falls back to binary serving on `UnicodeDecodeError`.
  - `serve_directory_listing()` — Generates an HTML directory listing; auto-redirects to `README.md` if present.
- `main()` — Entry point. Binds `HTTPServer` on the given port with graceful error on port conflicts.

The version string is maintained in two places: `viewmd.py:VERSION` and `pyproject.toml:[project].version`.

## Coding Guidelines

- **Max 200 lines per file.** If a file grows beyond this, split it into modules.
- **Max ~60 lines per function.** Extract helpers or break logic into smaller functions.
- Enforced by pylint (`max-module-lines=200`, `max-statements=30`), configured in `pyproject.toml`.

## Testing

Tests are in `tests/test_frontmatter.py` with fixtures in `tests/fixtures/`. Integration tests spin up a real HTTP server on a random port (`HTTPServer` bound to port 0) in a temp directory, copying fixture `.md` files there.
