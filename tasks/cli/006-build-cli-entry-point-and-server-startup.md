---
id: "006"
title: "Build CLI entry point and server startup"
status: pending
priority: high
dependencies: ["005"]
tags: ["go", "cli"]
created: 2026-03-07
---

# Build CLI entry point and server startup

## Objective

Wire up `cmd/viewmd/main.go` as the CLI entry point that parses args, starts the HTTP server, and handles graceful shutdown.

## Tasks

- [ ] Parse optional port argument from `os.Args` (default 8000)
- [ ] Bind HTTP server with the handler from task 005
- [ ] Detect port-in-use and print a clear error
- [ ] Print startup banner (version, URL, feature list) matching Python output
- [ ] Handle SIGINT/SIGTERM for graceful shutdown
- [ ] Define VERSION constant

## Acceptance Criteria

- `viewmd` starts on port 8000, `viewmd 3000` on port 3000
- Port conflict prints "Error: Port N is already in use." and exits
- Ctrl+C prints "Shutting down..." and exits cleanly
- Startup banner displays version, URL, and feature list
