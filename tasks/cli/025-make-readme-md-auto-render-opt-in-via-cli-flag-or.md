---
title: "Make README.md auto-render opt-in via CLI flag or config"
id: "025"
status: completed
priority: medium
type: feature
tags: ["config", "navigation"]
created: "2026-03-08"
---

# Make README.md auto-render opt-in via CLI flag or config

## Objective

Currently, navigating to a directory that contains a `README.md` automatically renders the README instead of showing the directory listing. This is unexpected — users should see the directory listing by default and choose to open the README themselves. Change the default to always show the directory listing, and add a CLI flag and/or config option to opt into the old auto-render behavior.

## Tasks

- [x] Remove the automatic README.md redirect from `serveDirectoryListing` / `ServeHTTP` in `handler.go`
- [x] Add a CLI flag (e.g. `--auto-readme`) to opt into auto-rendering README.md in directories
- [x] Wire the flag through to the handler (e.g. as a field on `Handler`)
- [x] Update existing tests that rely on the README auto-serve behavior
- [x] Add tests for both default (no auto-render) and opt-in (auto-render) modes

## Acceptance Criteria

- By default, navigating to a directory always shows the directory listing, even if it contains a README.md
- When the `--auto-readme` flag is passed (or equivalent config option is set), directories with a README.md render it automatically (current behavior)
- Existing directory listing features (sorting, parent link, etc.) are unaffected
- The flag is documented in `--help` output
