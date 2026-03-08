---
title: "Add inline image viewing"
id: "032"
status: in-progress
priority: medium
type: feature
tags: ["feature"]
touches: ["cli/handler", "cli/render"]
created: "2026-03-08"
---

# Add inline image viewing

## Objective

When a user clicks on an image file (PNG, JPG, GIF, SVG, WebP, etc.) in a directory listing, display it inline in the browser instead of showing "can't preview" or triggering a download. Images should render in a clean, centered view with the same page chrome (breadcrumbs, back navigation) as other file types.

## Tasks

- [ ] Add image file extension detection to the classify package
- [ ] Create an image rendering template/page
- [ ] Serve image files with correct Content-Type when accessed directly
- [ ] Route image file requests to the image viewer page
- [ ] Show image thumbnails or icons in directory listings (optional)
- [ ] Test with common formats: PNG, JPG, GIF, SVG, WebP

## Acceptance Criteria

- Clicking an image file in a directory listing displays it inline in the browser
- Image page includes breadcrumb navigation and back button
- Common image formats are supported (PNG, JPG, GIF, SVG, WebP)
- Raw image can still be accessed via `?raw=1`
- No layout issues with very large or very small images
