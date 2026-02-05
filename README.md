# viewmd

A simple HTTP server for viewing Markdown files in your browser.

## Installation

### Option 1: Using pipx (Recommended)

[pipx](https://pypa.github.io/pipx/) installs CLI tools in isolated environments, preventing dependency conflicts:

```bash
# Install pipx if you don't have it
python3 -m pip install --user pipx
python3 -m pipx ensurepath

# Install viewmd
pipx install .
```

### Option 2: Using pip

Install globally with pip:

```bash
pip install .
```

Or install in development mode (changes to source code take effect immediately):

```bash
pip install -e .
```

### Option 3: Direct symlink (Manual)

```bash
chmod +x viewmd.py
sudo ln -s "$(pwd)/viewmd.py" /usr/local/bin/viewmd
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
# If installed with pipx
pipx uninstall viewmd

# If installed with pip
pip uninstall viewmd
```

## Publishing to PyPI (Optional)

To make viewmd installable via `pip install viewmd` for everyone:

```bash
# Install build tools
pip install build twine

# Build the package
python -m build

# Upload to PyPI (requires account at pypi.org)
python -m twine upload dist/*
```

That's it. Simple.
