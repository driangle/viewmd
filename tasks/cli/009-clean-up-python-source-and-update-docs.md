---
id: "009"
title: "Archive Python source and update docs"
status: pending
priority: low
dependencies: ["008"]
tags: ["cleanup"]
touches: ["python", "root-docs"]
created: 2026-03-07
---

# Archive Python source and update docs

## Objective

Move the Python source into a `python/` archive directory (it's already published to PyPI and must be preserved) and update documentation to reflect the Go-based tool.

## Tasks

- [ ] Create `python/` directory at repo root
- [ ] Move `viewmd/` Python package into `python/viewmd/`
- [ ] Move `tests/` into `python/tests/`
- [ ] Move `pyproject.toml`, `requirements.txt`, `MANIFEST.in` into `python/`
- [ ] Remove `viewmd.egg-info/` and `dist/` (build artifacts, not source)
- [ ] Add `python/README.md` noting this is the archived PyPI-published Python version
- [ ] Update root README.md with Go install instructions (`go install` or binary download)
- [ ] Update CLAUDE.md with Go project description, commands, and architecture

## Acceptance Criteria

- Python source preserved under `python/` (viewmd package, tests, pyproject.toml)
- `python/` directory is self-contained (can still be built/tested from there)
- Build artifacts (`egg-info`, `dist`) removed (not archived)
- Root README reflects Go installation and usage
- CLAUDE.md accurately describes Go project structure
- `go build ./...` and `go test ./...` still pass
