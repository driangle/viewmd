---
id: "015"
title: "Port text, directory, and routing integration tests"
status: pending
priority: medium
tags: ["go", "testing"]
parent: "007"
dependencies: ["012", "013"]
created: 2026-03-07
---

# Port text, directory, and routing integration tests

## Objective

Port the Python `TestServeTextFile`, `TestServeDirectoryListing`, and `TestRouting` integration tests to Go.

## Tasks

- [ ] Test text file wrapped in `<pre>` tag (hello.py)
- [ ] Test text file shows filename in header
- [ ] Test text file content HTML-escaped
- [ ] Test dotfile served as text (.gitignore)
- [ ] Test known filename served as text (Makefile)
- [ ] Test binary content with text extension doesn't crash (broken.txt with invalid UTF-8)
- [ ] Test root lists files and subdirectories
- [ ] Test subdir with README.md renders markdown
- [ ] Test subdir without README.md shows directory listing
- [ ] Test subdir listing has parent link (..)
- [ ] Test markdown file routed to renderer
- [ ] Test binary file served raw (PNG)
- [ ] Test missing file returns 404
- [ ] Test URL-encoded paths work

## Acceptance Criteria

- All Python `TestServeTextFile` (6), `TestServeDirectoryListing` (4), and `TestRouting` (4) test cases pass
- Tests create appropriate temp files (hello.py, .gitignore, Makefile, image.png, broken.txt, docs/README.md, empty_dir/)
- Tests use real HTTP requests to a running server
