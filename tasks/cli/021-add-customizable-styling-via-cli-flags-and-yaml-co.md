---
title: "Add customizable styling via CLI flags and YAML config"
id: "021"
status: pending
priority: low
type: feature
tags: ["ui", "config"]
touches: ["cli/cmd", "cli/handler", "cli/render"]
created: "2026-03-07"
---

# Add customizable styling via CLI flags and YAML config

## Objective

Let users customize basic styling options for the rendered pages — such as font family, font size, max width, and theme/color scheme — via CLI flags or a YAML config file (e.g. `~/.config/viewmd/config.yaml` or `.viewmd.yaml` in the served directory).

## Tasks

- [ ] Define supported style options (e.g. `font-family`, `font-size`, `max-width`, `theme: light|dark`)
- [ ] Add YAML config file loading (check `.viewmd.yaml` in served directory, then `~/.config/viewmd/config.yaml`)
- [ ] Add CLI flags for each style option (e.g. `--font-size 16px`, `--theme dark`)
- [ ] CLI flags override config file values
- [ ] Inject style values into HTML templates as CSS custom properties or inline styles
- [ ] Use sensible defaults when no config or flags are provided
- [ ] Add tests for config loading, CLI flag parsing, and style injection

## Acceptance Criteria

- Users can create a `.viewmd.yaml` file to set styling preferences
- Users can pass `--font-size`, `--theme`, etc. as CLI flags
- CLI flags take precedence over config file values
- Config file values take precedence over defaults
- Style changes are reflected in the rendered HTML pages
- The server works correctly with no config file and no flags (defaults apply)
- Config file supports at minimum: `font-family`, `font-size`, `max-width`, `theme`
