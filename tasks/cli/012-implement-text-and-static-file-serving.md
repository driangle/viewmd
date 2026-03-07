---
id: "012"
title: "Implement text and static file serving"
status: completed
priority: high
tags: ["go", "http"]
parent: "005"
dependencies: ["003", "004"]
created: 2026-03-07
---

# Implement text and static file serving

## Objective

Implement the text file rendering path (HTML-escaped content in a pre block) and the static file passthrough for binary files.

## Tasks

- [x] Read text files and validate UTF-8 encoding
- [x] HTML-escape content and render using text page template
- [x] Fall back to binary serving (static file) when file contains invalid UTF-8
- [x] Serve unknown file types via `http.FileServer` passthrough (images, PDFs, etc.)
- [x] Handle errors with 500 responses

## Acceptance Criteria

- Text files served with HTML escaping in pre blocks
- Filename shown in header of text file pages
- Invalid UTF-8 files fall back to binary serving (not an error)
- Binary files served with correct content type
