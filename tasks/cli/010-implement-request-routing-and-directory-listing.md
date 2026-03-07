---
id: "010"
title: "Implement request routing and directory listing"
status: completed
priority: high
tags: ["go", "http"]
parent: "005"
dependencies: ["003", "004"]
created: 2026-03-07
---

# Implement request routing and directory listing

## Objective

Build the core HTTP handler in `apps/cli/internal/handler/` with the request routing logic and directory listing functionality, matching the Python `MarkdownHandler.do_GET()` dispatch.

## Tasks

- [x] Create `Handler` struct that implements `http.Handler`
- [x] Implement URL decoding and path resolution from request URL
- [x] Route root path (`/`) to directory listing
- [x] Route directory paths: serve README.md if present, otherwise show listing
- [x] Dispatch `.md`/`.markdown` to markdown serving (stub initially)
- [x] Dispatch text files to text serving (stub initially)
- [x] Dispatch other files to `http.FileServer` passthrough
- [x] Return 404 for missing files
- [x] Add request logging: `[HH:MM:SS] Request: /path`
- [x] Implement directory listing: sorted entries (dirs first), parent link for subdirs

## Acceptance Criteria

- Routing logic matches Python behavior for all path types
- Directory listing sorted with directories first, then files alphabetically
- Parent link (`..`) shown for subdirectories but not root
- 404 returned for non-existent paths
- Request logging printed to stdout with timestamp
