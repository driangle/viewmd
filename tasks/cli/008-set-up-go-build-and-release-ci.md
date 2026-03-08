---
id: "008"
title: "Set up Go build and release CI"
status: pending
priority: medium
dependencies: ["007"]
tags: ["go", "ci"]
touches: ["ci"]
created: 2026-03-07
---

# Set up Go build and release CI

## Objective

Replace the Python PyPI publishing CI with Go build/test/release workflows using GitHub Actions.

## Tasks

- [ ] Create GitHub Actions workflow for Go test/lint on push to main
- [ ] Create release workflow that cross-compiles binaries (linux/darwin/windows × amd64/arm64)
- [ ] Attach binaries to GitHub Releases
- [ ] Remove or replace the existing PyPI publish workflow

## Acceptance Criteria

- CI runs `go test ./...` and `go vet` on every push to main
- Release workflow triggered on GitHub Release creation
- Binaries produced for 6 platform/arch combinations
- Binaries attached to the GitHub Release automatically
