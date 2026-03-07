---
id: "015"
title: "Port text, directory, and routing integration tests"
status: completed
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

- [x] Test text file wrapped in `<pre>` tag (hello.py)
- [x] Test text file shows filename in header
- [x] Test text file content HTML-escaped
- [x] Test dotfile served as text (.gitignore)
- [x] Test known filename served as text (Makefile)
- [x] Test binary content with text extension doesn't crash (broken.txt with invalid UTF-8)
- [x] Test root lists files and subdirectories
- [x] Test subdir with README.md renders markdown
- [x] Test subdir without README.md shows directory listing
- [x] Test subdir listing has parent link (..)
- [x] Test markdown file routed to renderer
- [x] Test binary file served raw (PNG)
- [x] Test missing file returns 404
- [x] Test URL-encoded paths work

## Acceptance Criteria

- All Python `TestServeTextFile` (6), `TestServeDirectoryListing` (4), and `TestRouting` (4) test cases pass
- Tests create appropriate temp files (hello.py, .gitignore, Makefile, image.png, broken.txt, docs/README.md, empty_dir/)
- Tests use real HTTP requests to a running server
