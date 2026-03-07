---
id: "013"
title: "Set up Go integration test infrastructure"
status: completed
priority: medium
tags: ["go", "testing"]
parent: "007"
dependencies: ["010"]
created: 2026-03-07
---

# Set up Go integration test infrastructure

## Objective

Create the shared test helpers and fixture files needed by all integration test suites.

## Tasks

- [x] Copy fixture `.md` files (basic.md, plain.md, escaping.md, stripped.md) into `testdata/` directory
- [x] Create test helper function that starts a real HTTP server on `:0` (random port)
- [x] Helper should copy fixtures + additional files into a temp directory and serve from there
- [x] Helper should return base URL and cleanup function
- [x] Create a `get(baseURL, path)` helper that fetches a URL and returns the response body as string

## Acceptance Criteria

- `testdata/` contains all 4 fixture `.md` files from the Python tests
- Test helper starts a real server (not mocked) on a random available port
- Helper manages temp directory lifecycle (create + cleanup)
- `go test` passes with the infrastructure in place (even with no test cases yet)
