---
title: "Add code syntax highlighting for text/code files"
id: "031"
status: completed
priority: medium
type: feature
tags: ["feature"]
touches: ["cli/render"]
created: "2026-03-08"
---

# Add code syntax highlighting for text/code files

## Objective

Add language-aware syntax highlighting to text/code file viewing. Currently text files are displayed in monospace without color-coded syntax. Integrate a client-side syntax highlighting library (e.g. highlight.js or Prism) so code files are rendered with proper syntax colors.

## Tasks

- [x] Choose a syntax highlighting library (highlight.js or Prism)
- [x] Add the library CSS/JS to the text file template
- [x] Detect the language from the file extension and apply highlighting
- [x] Ensure highlighting works for all supported text file extensions
- [x] Support light and dark mode themes for syntax highlighting
- [x] Test with a variety of file types (Go, Python, JS, JSON, YAML, etc.)

## Acceptance Criteria

- Code files display with language-appropriate syntax coloring
- Language is auto-detected from file extension
- Syntax highlighting respects the current light/dark theme
- No JavaScript errors or layout regressions on non-code text files
