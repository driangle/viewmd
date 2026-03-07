---
id: "003"
title: "Implement file type classification"
status: pending
priority: high
dependencies: ["001"]
tags: ["go", "core"]
created: 2026-03-07
---

# Implement file type classification

## Objective

Implement `classify.IsTextFile(name string) bool` in `internal/classify/` matching the Python `MarkdownHandler.is_text_file()` logic.

## Tasks

- [ ] Define the TEXT_EXTENSIONS set (all 50+ extensions from Python)
- [ ] Define the TEXT_FILENAMES set (makefile, dockerfile, license, etc.)
- [ ] Implement dotfile detection (names starting with `.`)
- [ ] Write table-driven tests porting all `TestIsTextFile` cases

## Acceptance Criteria

- Same extensions/filenames recognized as the Python version
- Case-insensitive matching for extensions and filenames
- Dotfiles detected as text
- `.md`/`.markdown` NOT classified as text (handled separately)
- All 10 Python test cases pass
