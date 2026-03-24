---
title: "Export folder as ZIP with rendered HTML"
id: "036"
status: completed
priority: medium
type: feature
tags: ["export", "zip", "html"]
created: "2026-03-24"
---

# Export folder as ZIP with rendered HTML

## Objective

Add an "Export as ZIP (HTML)" action on directory pages. When triggered for a folder, the server generates a ZIP archive containing all files in that folder recursively. Markdown files are converted to rendered HTML before inclusion; all other files are included as-is.

## Tasks

- [x] Add a new HTTP endpoint (e.g. `GET /<path>?export=zip`) that accepts a directory path and produces a ZIP download
- [x] Implement recursive directory walking to collect all files under the target folder
- [x] For each `.md` file, render it to HTML using the existing Goldmark pipeline and include the `.html` output in the ZIP
- [x] For non-markdown files, include them in the ZIP unchanged, preserving relative paths
- [x] Set appropriate response headers (`Content-Type: application/zip`, `Content-Disposition: attachment`)
- [x] Add an "Export as ZIP" button/link to the directory listing template
- [x] Add unit tests for ZIP generation (correct file list, markdown converted, paths preserved)
- [x] Add integration test using httptest with a temp directory fixture
- [x] Update documentation (CLI reference, README) to describe the export feature

## Acceptance Criteria

- Clicking "Export as ZIP" on a directory page downloads a `.zip` file
- The ZIP contains all files from the directory recursively, preserving relative folder structure
- All `.md` files in the ZIP are converted to rendered HTML (`.html` extension)
- Non-markdown files are included verbatim with no modification
- The ZIP file name reflects the exported folder name (e.g. `my-folder.zip`)
- Empty subdirectories are not included
- Works correctly for nested directories with mixed file types
