---
title: "Add full-text search with keyboard-driven UI"
id: "030"
status: in-progress
priority: high
type: feature
tags: ["search", "ui", "keyboard"]
created: "2026-03-08"
---

# Add full-text search with keyboard-driven UI

## Objective

Add a search feature to the viewmd UI that lets users search the entire directory tree by file name, file contents, or both. The search should activate automatically when the user starts typing, filtering results in real time. The entire flow must be fully keyboard-navigable — no mouse required.

## Tasks

- [ ] Add a search input bar to the directory template (visible on all directory pages)
- [ ] Implement a backend search API endpoint (`/search?q=...&mode=name|content|both`) that recursively searches the served directory tree
  - [ ] File name matching (case-insensitive substring/fuzzy)
  - [ ] File content matching (case-insensitive substring grep)
  - [ ] Combined mode (search both name and content)
- [ ] Build a keyboard-driven search results UI
  - [ ] Auto-activate search mode when user starts typing (or on `/` keypress)
  - [ ] Display results in a filterable list as the user types (debounced requests)
  - [ ] Allow switching between search modes (name / content / both) via keyboard shortcut (e.g., Tab or Ctrl+M)
  - [ ] Navigate results with arrow keys (Up/Down) and open selected result with Enter
  - [ ] Dismiss search with Escape, returning focus to the directory listing
- [ ] Show result context: for content matches, display the matching line or snippet
- [ ] Handle edge cases: empty query clears results, no results found message, large directories (limit results or paginate)

## Acceptance Criteria

- Typing in a directory view immediately enters search mode and filters results
- Users can search by file name only, file contents only, or both
- Search mode is clearly indicated and switchable via keyboard
- Results list is navigable with arrow keys; Enter opens the selected file
- Escape exits search and returns to normal directory view
- No mouse interaction is required for the full search workflow
- Search works across the entire directory tree (recursive)
- Content match results show a snippet of the matching line
