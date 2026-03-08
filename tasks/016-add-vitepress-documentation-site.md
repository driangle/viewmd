---
title: "Add VitePress documentation site"
id: "016"
status: completed
priority: medium
type: feature
tags: ["docs", "vitepress"]
touches: ["docs"]
created: "2026-03-07"
---

# Add VitePress documentation site

## Objective

Add a VitePress-powered documentation site for viewmd, similar to the [taskmd docs site](https://driangle.github.io/taskmd/). The site should document both the Python CLI (legacy) and the Go CLI (in development), and be deployed to GitHub Pages.

## Tasks

- [x] Initialize VitePress project in a `docs/` directory (`npm init` + `vitepress` dependency)
- [x] Configure VitePress (`docs/.vitepress/config.ts`) with nav, sidebar, and site metadata
- [x] Create homepage (`docs/index.md`) with hero section and feature highlights
- [x] Write "Getting Started" section: installation, quick start, basic usage
- [x] Write "Guide" section: CLI usage, markdown rendering features, frontmatter support, directory listing
- [x] Write "Reference" section: CLI flags, configuration, supported file types
- [x] Add GitHub Pages deployment workflow (`.github/workflows/docs.yml`)
- [x] Add `docs:dev` and `docs:build` scripts to root or docs-level `package.json`
- [x] Verify the site builds and deploys correctly

## Acceptance Criteria

- `docs/` directory contains a working VitePress project
- `npm run docs:dev` starts a local dev server with the documentation site
- `npm run docs:build` produces a static site without errors
- Site has a homepage with hero, feature cards, and navigation
- Sidebar includes Getting Started, Guide, and Reference sections
- GitHub Actions workflow deploys the site to GitHub Pages on push to `main`
- Site structure and styling follows the same pattern as [driangle.github.io/taskmd](https://driangle.github.io/taskmd/)
