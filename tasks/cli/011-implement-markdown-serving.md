---
id: "011"
title: "Implement markdown serving with goldmark"
status: pending
priority: high
tags: ["go", "http"]
parent: "005"
dependencies: ["002", "004"]
created: 2026-03-07
---

# Implement markdown serving with goldmark

## Objective

Implement the markdown file rendering path: read `.md` files, parse frontmatter, convert markdown to HTML using goldmark, and serve the rendered page.

## Tasks

- [ ] Configure goldmark with fenced code blocks and tables extensions
- [ ] Read markdown file, parse frontmatter using `apps/cli/internal/frontmatter`
- [ ] Convert markdown body to HTML
- [ ] Compute base URL from file's parent directory
- [ ] Render full HTML page using `apps/cli/internal/render` templates
- [ ] Handle errors (file read failures, rendering errors) with 500 responses

## Acceptance Criteria

- Markdown rendered with fenced code block and table support
- Frontmatter stripped from body and rendered as HTML table above content
- Base href set correctly for relative links
- Errors return 500 with descriptive message
