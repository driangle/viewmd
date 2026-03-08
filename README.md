# viewmd

A minimal file viewer for your browser. Point it at a directory and browse everything — markdown rendered beautifully, code syntax-highlighted, images inline.

**View everything, render what you can, manage nothing.**

**[Documentation](https://driangle.github.io/viewmd/)**

## Installation

```bash
brew install driangle/tap/viewmd
```

Or from source: `go install github.com/driangle/viewmd/apps/cli/cmd/viewmd@latest`

## Quick Start

```bash
viewmd          # Serves current directory on port 8000
viewmd 3000     # Custom port
viewmd -a       # Show all files, not just markdown
viewmd -w       # Watch for changes, auto-reload browser
```

## Screenshots

**Markdown rendering** — `.md` files are rendered as styled HTML with full typography support.

![Markdown rendering](docs/public/images/viewmd_md.png?v=2)

**Directory browsing** — Navigate your project with a clean file listing, keyboard shortcuts, and search.

![Directory browsing](docs/public/images/viewmd_dir.png?v=2)

**Code and text files** — View any text file with syntax highlighting, copy support, and dark mode.

![Code and text file viewing](docs/public/images/viewmd_license.png?v=2)

## What It Does

- **Markdown** — Rendered HTML with styled typography
- **Code/text** — Syntax-highlighted display
- **Images** — Inline rendering
- **Directories** — File listing with keyboard navigation and search
- **Other files** — Download or "can't preview" message
