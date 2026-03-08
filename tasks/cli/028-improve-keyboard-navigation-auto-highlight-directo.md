---
title: "Improve keyboard navigation: auto-highlight directory on enter and back"
id: "028"
status: completed
priority: medium
type: feature
tags: ["keyboard", "navigation"]
created: "2026-03-08"
---

# Improve keyboard navigation: auto-highlight directory on enter and back

## Objective

Improve keyboard navigation in the directory listing so that navigating between directories feels seamless. When entering a directory, automatically highlight the first item. When navigating back to a parent directory, automatically highlight the directory the user just came from.

## Tasks

- [x] Pass the "came from" directory name via URL fragment or query param when navigating back (e.g., `?from=subdir`)
- [x] On directory page load, check for a `from` parameter and find the matching entry in `#file-list`
- [x] If a `from` match is found, auto-highlight that list item; otherwise, highlight the first item
- [x] When entering a directory (Enter/ArrowRight), highlight the first item on the new page
- [x] Ensure parent `..` link also passes `from` context when using ArrowLeft/Backspace navigation

## Acceptance Criteria

- Entering a directory via keyboard highlights the first list item automatically
- Pressing ArrowLeft/Backspace to go to the parent directory highlights the directory the user just left
- Clicking links (mouse navigation) does not auto-highlight anything (current behavior preserved)
- All existing keyboard shortcuts (ArrowUp, ArrowDown, Enter, ArrowRight, ArrowLeft, Backspace) continue to work correctly
- No regressions in existing Go tests (`make test` passes)
