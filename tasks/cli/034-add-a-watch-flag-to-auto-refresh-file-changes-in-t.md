---
title: "Add a --watch flag to auto-refresh file changes in the web viewer"
id: "034"
status: in-progress
priority: medium
type: feature
tags: []
created: "2026-03-08"
---

# Add a --watch flag to auto-refresh file changes in the web viewer

## Objective

Add a `--watch` mode to the `viewmd` CLI that monitors the served path for file and directory changes and ensures the browser view reflects updates automatically without manually restarting the server.

## Tasks

- [ ] Define CLI behavior for `--watch`, including defaults, supported scope (single file vs. directory), and expected UX when changes are detected.
- [ ] Implement filesystem watching in the CLI/server path startup flow using a reliable watcher strategy for nested directories and common editor save patterns.
- [ ] Update HTTP rendering/asset invalidation behavior so modified Markdown, text/code, images, and directory listings are reflected on refresh immediately after changes.
- [ ] Add automatic browser refresh signaling in watch mode (for example via SSE/WebSocket/polling endpoint) so open viewer tabs update when watched files change.
- [ ] Add tests for watch mode behavior, including change detection and updated content visibility, with temporary directories and test fixtures.
- [ ] Document `--watch` usage and caveats in CLI help output and relevant project docs.

## Acceptance Criteria

- Running `viewmd --watch <path>` starts the server and keeps content in sync with local file changes without requiring process restart.
- Editing a Markdown file in the watched scope causes the open viewer page to show updated rendered output automatically.
- Editing non-Markdown previewable files (text/code/images) in the watched scope is reflected in the viewer after change notification.
- Directory view content (file list and README display) updates when files are added, removed, or renamed in watched directories.
- Watch mode behavior is covered by automated tests for at least one file-update path and one directory-change path.
- CLI help and docs clearly describe `--watch`, its behavior, and limitations.
