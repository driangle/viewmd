---
title: "Improve UI design — minimal and functional refresh"
id: "027"
status: completed
priority: medium
type: feature
tags: ["ui", "frontend"]
created: "2026-03-08"
---

# Improve UI design — minimal and functional refresh

## Objective

Refresh the Go CLI's HTML templates (`apps/cli/internal/render/templates/`) with a cleaner, more polished visual design while keeping the UI minimal and functional. Use the `/frontend-design` skill to produce production-grade styles. The result should feel purposeful and modern — not flashy.

## Tasks

- [ ] Audit current templates (`markdown.html`, `text.html`, `directory.html`) and identify visual inconsistencies and rough spots
- [ ] Redesign the shared color palette, typography, and spacing for a cohesive look across all three page types
- [ ] Improve the directory listing page — better file/folder icons, clearer visual hierarchy, refined keyboard-active highlight
- [ ] Improve the markdown page — better frontmatter table styling, improved code block and blockquote rendering, refined navigation bar
- [ ] Improve the text file page — align card styling with the rest, better header/copy button treatment
- [ ] Ensure all pages remain responsive and usable on small screens
- [ ] Verify existing keyboard navigation and copy-to-clipboard still work after changes
- [ ] Run `make test` to confirm nothing is broken

## Acceptance Criteria

- All three Go templates share a consistent visual language (colors, type, spacing)
- Pages look noticeably cleaner and more polished than before
- No new JS dependencies — CSS-only improvements where possible
- Existing functionality (keyboard nav, copy button, parent nav) is preserved
- `make test` passes
