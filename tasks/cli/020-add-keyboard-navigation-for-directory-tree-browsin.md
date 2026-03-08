---
title: "Add keyboard navigation for directory tree browsing"
id: "020"
status: completed
priority: low
type: feature
tags: ["ui", "navigation"]
touches: ["cli/render"]
created: "2026-03-07"
---

# Add keyboard navigation for directory tree browsing

## Objective

Add simple keyboard shortcuts to navigate the directory tree in the browser. Users should be able to move up/down through file listings, enter files/folders, and go back to the parent directory — all without touching the mouse.

## Tasks

- [x] Add `keydown` event listener to directory listing pages for arrow key navigation (up/down to highlight items)
- [x] Implement Enter key to open the highlighted file or folder
- [x] Implement Backspace or left arrow to navigate up to the parent directory
- [x] Add visual highlight/focus indicator on the currently selected list item
- [x] On file view pages, support Backspace or left arrow to return to the parent directory listing
- [x] Add a small keyboard shortcut hint (e.g. "↑↓ navigate · Enter open · ← back") at the bottom of directory pages
- [x] Add tests for keyboard event handling

## Acceptance Criteria

- On directory listing pages, up/down arrow keys move a visible highlight through the file list
- Pressing Enter on a highlighted item navigates to that file or folder
- Pressing Backspace or left arrow navigates to the parent directory (from both directory and file views)
- The currently highlighted item has a clear visual indicator
- Keyboard navigation does not interfere with normal browser shortcuts (e.g. scrolling, browser back)
- Navigation works without JavaScript errors in modern browsers
