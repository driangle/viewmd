---
title: "Show file path breadcrumb when viewing a file"
id: "029"
status: pending
priority: medium
type: feature
tags: ["ui", "navigation"]
created: "2026-03-08"
---

# Show file path breadcrumb when viewing a file

## Objective

Display the current file's path when viewing markdown or text files, so the user always knows where they are in the directory tree. Currently only the `..` parent link is shown, with no indication of the file's location.

## Tasks

- [ ] Add a file path display to the markdown template (`markdown.html`) showing the path relative to the serve root
- [ ] Add a file path display to the text template (`text.html`) in a similar location
- [ ] Pass the display path from the handler to the template data (may need a new template field)
- [ ] Style the path as a breadcrumb with clickable path segments linking to parent directories
- [ ] Ensure the breadcrumb replaces or integrates with the existing `..` parent-nav link

## Acceptance Criteria

- The current file's relative path is visible on both markdown and text file views
- Each path segment is clickable and navigates to that directory
- The breadcrumb styling is consistent with the design system (uses CSS custom properties)
- Existing keyboard navigation (ArrowLeft/Backspace to go back) continues to work
- All existing Go tests pass (`make test`)
