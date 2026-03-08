---
title: "Add light/dark mode toggle with localStorage persistence"
id: "022"
status: completed
priority: low
type: feature
tags: ["ui", "theme"]
touches: ["cli/render"]
created: "2026-03-07"
---

# Add light/dark mode toggle with localStorage persistence

## Objective

Add a light/dark mode toggle to all viewmd pages. The user's preference should be saved in `localStorage` so it persists across page loads and navigation. On first visit, default to the system preference (`prefers-color-scheme`).

## Tasks

- [ ] Define light and dark color schemes using CSS custom properties (background, text, links, code blocks, borders)
- [ ] Add a toggle button (e.g. sun/moon icon) to the page layout, consistent across all page types
- [ ] Implement JavaScript to toggle a class/attribute on `<html>` (e.g. `data-theme="dark"`)
- [ ] Save the selected theme to `localStorage` on toggle
- [ ] On page load, read `localStorage` first; fall back to `prefers-color-scheme` media query
- [ ] Apply the theme early (inline `<script>` in `<head>`) to prevent flash of wrong theme
- [ ] Update all templates (markdown, text, directory listing) to include the toggle and theme styles

## Acceptance Criteria

- A visible toggle button switches between light and dark mode on all page types
- The selected theme persists across page reloads and navigation via `localStorage`
- First-time visitors see a theme matching their OS/browser preference
- No flash of unstyled/wrong-theme content on page load
- Both themes have readable contrast for all content (text, code, links, tables)
- The toggle is unobtrusive and consistently positioned across pages
