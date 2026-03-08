---
title: "Add support for configurable ignore patterns"
id: "033"
status: completed
priority: medium
type: feature
tags: []
created: "2026-03-08"
---

# Add support for configurable ignore patterns

## Objective

Allow users to configure which folders or files viewmd should ignore when listing directory contents. This enables hiding irrelevant entries like `node_modules`, `.git`, build artifacts, or any user-specified patterns from directory views.

## Tasks

- [x] Define a configuration format for ignore patterns (e.g. `.viewmdignore` file or CLI flags)
- [x] Implement pattern matching logic (support glob patterns like `*.log`, `node_modules/`, `.git/`)
- [x] Integrate ignore filtering into the directory listing handler
- [x] Add built-in default ignore patterns (e.g. `.git`) with ability to override
- [x] Add CLI flag to specify additional ignore patterns (e.g. `--ignore`)
- [x] Add unit tests for pattern matching and filtering logic
- [x] Add integration test for directory listing with ignore patterns

## Acceptance Criteria

- Users can create a `.viewmdignore` file with glob patterns to hide files/folders from directory listings
- A `--ignore` CLI flag allows specifying additional patterns at launch
- Default patterns (`.git`) are applied unless explicitly overridden
- Ignored files/folders do not appear in directory listing pages
- Pattern matching supports standard glob syntax (wildcards, directory patterns)
- Existing behavior is unchanged when no ignore configuration is present
