---
title: "Add copy-to-clipboard button for page contents"
id: "018"
status: completed
priority: low
type: feature
tags: ["ui"]
touches: ["cli/render"]
created: "2026-03-07"
---

# Add copy-to-clipboard button for page contents

## Objective

Add a "Copy" button to each file view page (markdown and text) that copies the raw file contents to the clipboard. This lets users quickly grab the full source content without manually selecting text.

## Tasks

- [x] Add a "Copy" button element to the markdown file view template
- [x] Add a "Copy" button element to the text file view template
- [x] Implement clipboard copy using the Clipboard API (`navigator.clipboard.writeText`)
- [x] Copy the raw/original file content (not the rendered HTML)
- [x] Add visual feedback on copy (e.g. button text briefly changes to "Copied!")
- [x] Style the button to be unobtrusive (e.g. top-right corner, small icon/text)

## Acceptance Criteria

- A "Copy" button is visible on markdown and text file view pages
- Clicking the button copies the raw file contents (original markdown/text source) to the clipboard
- The button provides visual feedback indicating the copy succeeded
- The button is positioned unobtrusively and does not interfere with reading
- Works in modern browsers that support the Clipboard API
