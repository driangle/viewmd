---
id: "005"
title: "Implement HTTP handler and routing"
status: completed
priority: high
dependencies: ["002", "003", "004"]
tags: ["go", "http"]
created: 2026-03-07
---

# Implement HTTP handler and routing

## Objective

Implement the HTTP request handler in `apps/cli/internal/handler/` that routes requests to the correct rendering path, matching the Python `MarkdownHandler.do_GET()` behavior.

## Tasks

- [x] Implement URL decoding and path resolution
- [x] Route root path to directory listing
- [x] Route directories with README.md to markdown rendering
- [x] Route directories without README.md to directory listing
- [x] Route `.md`/`.markdown` files to markdown rendering (with goldmark conversion using fenced_code, tables extensions)
- [x] Route text files to text rendering (with binary fallback on invalid UTF-8)
- [x] Route other files to `http.FileServer` passthrough
- [x] Return 404 for missing files
- [x] Add request logging with timestamp `[HH:MM:SS] Request: /path`

## Acceptance Criteria

- All routing paths match Python behavior
- Markdown rendered with fenced code blocks and tables support
- Text files HTML-escaped and served in pre blocks
- Invalid UTF-8 text files fall back to binary serving
- Binary files (images, etc.) served as-is
- 404 returned for non-existent files

## Sub-tasks

- **010** — Implement request routing and directory listing
- **011** — Implement markdown serving with goldmark
- **012** — Implement text and static file serving
