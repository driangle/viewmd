# viewmd

A minimal file viewer for your browser. Point it at a directory and browse everything — markdown rendered beautifully, code syntax-highlighted, images inline.

**View everything, render what you can, manage nothing.**

## Installation

### From source

```bash
go install github.com/user/viewmd/apps/cli/cmd/viewmd@latest
```

### Build locally

```bash
make build
# Binary at apps/cli/bin/viewmd
```

## Quick Start

```bash
viewmd          # Serves current directory on port 8000
viewmd 3000     # Custom port
```

## What It Does

- **Markdown files** (`.md`) — Rendered as HTML with nice styling
- **Code/text files** — Syntax-highlighted display
- **Images** — Inline rendering
- **Directories** — File listing with README auto-display
- **Other files** — Download or "can't preview" message

## Python version

The original Python version is archived in [`python/`](./python/) and available on [PyPI](https://pypi.org/project/viewmd/) (`pipx install viewmd`).
