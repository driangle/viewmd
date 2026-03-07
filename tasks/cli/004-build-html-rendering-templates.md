---
id: "004"
title: "Build HTML rendering templates"
status: pending
priority: high
dependencies: ["001"]
tags: ["go", "templates"]
created: 2026-03-07
---

# Build HTML rendering templates

## Objective

Create Go `html/template` templates in `internal/render/` that produce the same HTML/CSS output as the Python `render.py` functions, embedded via `embed.FS`.

## Tasks

- [ ] Create markdown page template (frontmatter table, base href, body HTML, full CSS)
- [ ] Create text file page template (filename header, pre block, CSS)
- [ ] Create directory listing template (sorted entries, parent link, CSS)
- [ ] Implement render functions: `RenderMarkdownPage()`, `RenderTextPage()`, `RenderDirectoryPage()`
- [ ] Embed templates using `//go:embed` directive
- [ ] Write unit tests verifying HTML output structure

## Acceptance Criteria

- CSS matches the Python version's styling
- Frontmatter renders as an HTML table with `.frontmatter` and `.fm-key` classes
- HTML auto-escaping via `html/template` prevents XSS
- Templates embedded in binary (no external files needed)
