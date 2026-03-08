---
title: "Support global daemon mode with multi-folder serving"
id: "019"
status: pending
priority: medium
type: feature
tags: ["cli", "server"]
touches: ["cli/cmd", "cli/handler"]
created: "2026-03-07"
---

# Support global daemon mode with multi-folder serving

## Objective

Allow viewmd to run as a single long-lived process (globally or in the background) and let users dynamically add folders from anywhere on their computer. Instead of starting a new server per directory, a single viewmd instance serves multiple registered folders, each accessible via its own route (e.g. `/projects/`, `/notes/`).

## Tasks

- [ ] Design the multi-folder routing scheme (e.g. `/<folder-name>/` prefix per registered folder)
- [ ] Add a `viewmd add <path>` command that registers a folder with a running instance
- [ ] Add a `viewmd remove <path>` command to unregister a folder
- [ ] Add a `viewmd list` command to show currently served folders
- [ ] Implement a mechanism for the running server to discover new folders (e.g. config file, Unix socket, or HTTP API)
- [ ] Add a root index page that lists all registered folders with links
- [ ] Add a `viewmd start` command for daemon/background mode (or reuse existing default)
- [ ] Persist registered folders across server restarts (e.g. config file in `~/.config/viewmd/`)
- [ ] Add tests for multi-folder routing and folder registration

## Acceptance Criteria

- A single viewmd process can serve files from multiple directories simultaneously
- Users can run `viewmd add ~/Documents/notes` (or similar) from any directory to register a folder
- The root URL (`/`) lists all registered folders
- Each folder is accessible under its own URL prefix
- Registered folders persist across server restarts
- `viewmd list` shows all currently registered folders
- `viewmd remove <path>` stops serving a previously added folder
