# Quick Start

## Basic Usage

Navigate to any directory and run:

```bash
viewmd
```

This starts a local server on port 8000. Open `http://localhost:8000` in your browser to browse the directory.

## Custom Port

```bash
viewmd 3000
```

## Show All Files

By default, viewmd only shows markdown files and directories containing markdown. To show all files:

```bash
viewmd -a
```

## Auto-display README

Automatically render a directory's README.md when you navigate to it:

```bash
viewmd --auto-readme
```

## Watch Mode

Auto-reload the browser when files change — useful for live editing:

```bash
viewmd --watch
```

## Combine Flags

```bash
viewmd -a -w --auto-readme 3000
```

This shows all files, watches for changes, auto-displays READMEs, and serves on port 3000.
