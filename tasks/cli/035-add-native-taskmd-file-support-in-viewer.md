---
title: "Add native taskmd file support in viewer"
id: "035"
status: pending
priority: medium
type: feature
effort: large
tags: ["taskmd", "integration"]
created: "2026-03-08"
---

# Add native taskmd file support in viewer

## Objective

Add native support for [taskmd](https://github.com/driangle/taskmd) markdown files in viewmd. When a `.md` file contains valid taskmd frontmatter (has `id`, `title`, and `status` fields), viewmd should detect it and render the frontmatter metadata in a structured, visually distinct way — showing status badges, priority indicators, tags, and other task fields above the rendered markdown body. In directory listings, taskmd files should show inline metadata (status, priority) alongside the filename.

This feature should be opt-in via `.viewmd.yaml` configuration so users who don't use taskmd aren't affected.

## Tasks

- [ ] Add `taskmd` section to `.viewmd.yaml` config (e.g., `taskmd: { enabled: true }`)
- [ ] Add the `github.com/driangle/taskmd/sdk/go` module as a dependency
- [ ] Create `internal/taskmd/` package to detect and parse taskmd files using the Go SDK's `parser.ParseTaskContent()`
- [ ] Detect taskmd files by checking for `id`, `title`, and `status` fields in YAML frontmatter
- [ ] Render taskmd frontmatter as a structured metadata block on file view pages (status badge, priority, effort, tags, dependencies, owner, dates)
- [ ] Add CSS styling for task metadata (status colors, priority indicators, tag chips)
- [ ] In directory view, show inline task metadata (status icon/badge, priority) next to taskmd files
- [ ] Strip taskmd frontmatter from the rendered markdown body (avoid duplicate display)
- [ ] Add unit tests for taskmd detection and metadata extraction
- [ ] Add integration tests for directory and file view rendering with taskmd files
- [ ] Update documentation (README, CLI reference) to describe the taskmd integration and config option

## Acceptance Criteria

- When `taskmd.enabled: true` is set in `.viewmd.yaml`, taskmd files are detected and rendered with structured metadata
- When `taskmd.enabled` is `false` or absent, taskmd files render as regular markdown (no special treatment)
- The file view page shows a metadata header with status, priority, effort, tags, owner, dependencies, and dates — styled with appropriate colors and badges
- Directory listings show inline status and priority indicators next to taskmd file entries
- The taskmd Go SDK (`github.com/driangle/taskmd/sdk/go/parser`) is used for parsing — no custom frontmatter parsing for task fields
- All existing markdown rendering continues to work unchanged
- Unit tests cover taskmd detection, metadata extraction, and edge cases (non-taskmd files, partial frontmatter)
