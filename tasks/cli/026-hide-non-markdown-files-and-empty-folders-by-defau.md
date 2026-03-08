---
title: "Hide non-markdown files and empty folders by default with opt-out flag"
id: "026"
status: completed
priority: medium
type: feature
tags: []
touches: ["cli/cmd", "cli/handler", "cli/render"]
created: "2026-03-08"
---

# Hide non-markdown files and empty folders by default with opt-out flag

## Objective

By default, directory listings should only show Markdown files (`.md`, `.markdown`) and folders that contain at least one Markdown file (recursively). Non-markdown files and empty/non-markdown-only folders should be hidden. Users can opt out of this behavior via a CLI flag or YAML config to show all files as before.

## Tasks

- [x] Add a `--show-all` / `-a` CLI flag that disables markdown-only filtering (default: off, meaning filtering is on)
- [x] Support a `show_all_files` option in a YAML config file (e.g., `.viewmd.yaml` or similar)
- [x] Update directory listing logic to filter out non-markdown files when filtering is enabled
- [x] Update directory listing logic to hide folders that contain no markdown files (recursively) when filtering is enabled
- [x] Ensure README.md auto-redirect still works when filtering is enabled
- [x] Add unit tests for the filtering logic
- [x] Add integration tests for CLI flag and config-based toggling

## Acceptance Criteria

- Directory listings hide non-markdown files by default
- Folders with no markdown files (at any depth) are hidden from listings by default
- Passing `--show-all` or `-a` shows all files and folders (original behavior)
- Setting `show_all_files: true` in YAML config shows all files and folders
- CLI flag takes precedence over YAML config
- README.md auto-redirect continues to work correctly
- Existing tests still pass
