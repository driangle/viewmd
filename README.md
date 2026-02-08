# viewmd

A simple HTTP server for viewing Markdown files in your browser.

## Installation

```bash
pipx install viewmd
```

## Quick Start

```bash
# Run from any directory
viewmd

# Open http://localhost:8000
```

## What It Does

- **Markdown files** (`.md`) - Rendered as HTML with nice styling
- **Text files** (`.py`, `.json`, `.gitignore`, etc.) - Displayed in browser
- **Directories** - Shows file listing, auto-displays `README.md`
- **Other files** - Served normally (images, PDFs, etc.)

## Usage

```bash
viewmd          # Starts on port 8000
viewmd 3000     # Custom port
```

## Uninstallation

```bash
pipx uninstall viewmd
```
