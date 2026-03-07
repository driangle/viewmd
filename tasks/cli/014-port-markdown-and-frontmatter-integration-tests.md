---
id: "014"
title: "Port markdown and frontmatter integration tests"
status: completed
priority: medium
tags: ["go", "testing"]
parent: "007"
dependencies: ["011", "013"]
created: 2026-03-07
---

# Port markdown and frontmatter integration tests

## Objective

Port the Python `TestServeMarkdownFrontmatter` integration tests to Go, verifying markdown rendering and frontmatter display via real HTTP requests.

## Tasks

- [x] Test frontmatter rendered as HTML table (basic.md: div.frontmatter, table, fm-key class, title/date values)
- [x] Test no frontmatter produces no table (plain.md)
- [x] Test frontmatter HTML escaping (escaping.md: no raw `<script>` tags)
- [x] Test frontmatter stripped from body (stripped.md: "status: draft" not in body after frontmatter div)
- [x] Test frontmatter CSS classes present in output

## Acceptance Criteria

- All 5 Python `TestServeMarkdownFrontmatter` test cases pass as Go tests
- Tests use real HTTP requests to a running server
