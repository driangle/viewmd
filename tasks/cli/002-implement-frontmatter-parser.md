---
id: "002"
title: "Implement frontmatter parser"
status: completed
priority: high
dependencies: ["001"]
tags: ["go", "core"]
created: 2026-03-07
---

# Implement frontmatter parser

## Objective

Implement `frontmatter.Parse(content string) (map[string]string, string)` in `apps/cli/internal/frontmatter/` that exactly matches the Python `parse_frontmatter()` behavior.

## Tasks

- [ ] Implement `Parse()` function: detect `---` delimiters, extract key:value pairs, return remaining body
- [ ] Handle edge cases: no frontmatter, single delimiter, empty values, colons in values, lines without colons, empty frontmatter block
- [ ] Return `nil` map (not empty) when no frontmatter is found
- [ ] Write table-driven tests porting all `TestParseFrontmatter` cases from Python

## Acceptance Criteria

- All 12 Python frontmatter test cases pass as Go tests
- `nil` returned for content without frontmatter
- Values with colons (URLs, timestamps) parsed correctly
- Whitespace trimmed from keys and values
