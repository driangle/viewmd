---
id: "006"
title: "Build CLI entry point and server startup"
status: completed
priority: high
dependencies: ["005"]
tags: ["go", "cli"]
created: 2026-03-07
---

# Build CLI entry point and server startup

## Objective

Wire up `apps/cli/cmd/viewmd/main.go` as the CLI entry point that parses args, starts the HTTP server, and handles graceful shutdown.

## Tasks

- [x] Parse optional port argument from `os.Args` (default 8000)
- [x] Bind HTTP server with the handler from task 005
- [x] Detect port-in-use and print a clear error
- [x] Print startup banner (version, URL, feature list) matching Python output
- [x] Handle SIGINT/SIGTERM for graceful shutdown
- [x] Define VERSION constant

## Acceptance Criteria

- `viewmd` starts on port 8000, `viewmd 3000` on port 3000
- Port conflict prints "Error: Port N is already in use." and exits
- Ctrl+C prints "Shutting down..." and exits cleanly
- Startup banner displays version, URL, and feature list
