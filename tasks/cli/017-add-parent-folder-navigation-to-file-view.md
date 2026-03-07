---
title: "Add parent folder navigation to file view"
id: "017"
status: pending
priority: low
type: feature
tags: ["ui", "navigation"]
created: "2026-03-07"
---

# Add parent folder navigation to file view

## Objective

Add lightweight navigation to the file view so users can easily navigate up to the containing folder. When viewing a rendered markdown or text file, there should be a simple link (e.g. ".. parent folder" or a breadcrumb) that takes the user back to the directory listing of the parent folder.

## Tasks

- [ ] Add a "navigate up" link/element to the markdown file view template
- [ ] Add the same link to the text file view template
- [ ] Ensure the link resolves to the correct parent directory path
- [ ] Style the link to be minimal and unobtrusive (top of page, small text)
- [ ] Add tests for the parent navigation link in rendered output

## Acceptance Criteria

- When viewing any file (markdown or text), a link to the parent directory is visible near the top of the page
- Clicking the link navigates to the directory listing of the containing folder
- The link works correctly for files at any depth (e.g. `/foo/bar/baz.md` links to `/foo/bar/`)
- The navigation element is lightweight and does not clutter the file view
- Files served from the root directory link to `/`
