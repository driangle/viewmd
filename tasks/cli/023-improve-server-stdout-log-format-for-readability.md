---
title: "Improve server stdout log format for readability"
id: "023"
status: pending
priority: low
type: chore
tags: ["logging", "dx"]
touches: ["cli/cmd", "cli/handler"]
created: "2026-03-07"
---

# Improve server stdout log format for readability

## Description

Review and refine the stdout log output of the viewmd server. Logs should be concise, minimal, aesthetically clean, and follow conventional formats (similar to tools like Vite, Hugo, or Caddy). Avoid noisy or verbose output — every line should be useful at a glance.

## Tasks

- [ ] Audit current log output from both the Python and Go servers
- [ ] Define a consistent log format (e.g. `HH:MM:SS METHOD /path STATUS duration`)
- [ ] Keep startup message minimal: address, port, and served directory on one or two lines
- [ ] Use short, human-readable timestamps (e.g. `15:04:05` not full ISO)
- [ ] Color-code or visually distinguish status codes if terminal supports it (200 green, 404 yellow, 500 red)
- [ ] Suppress redundant info (e.g. repeated host, user-agent, HTTP version)
- [ ] Ensure quiet/clean shutdown message
- [ ] Update both Python and Go implementations to match the agreed format

## Acceptance Criteria

- Server startup prints a clean, minimal banner with address and served directory
- Each request logs on a single line: timestamp, method, path, status, and duration
- Timestamps are short and human-readable (not full ISO-8601)
- No redundant or noisy fields in the default log output
- Log format is consistent between Python and Go implementations
- Output is easy to scan when tailing in a terminal
