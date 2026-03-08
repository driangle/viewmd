# viewmd (Python — Archived)

This is the original Python implementation of viewmd, published on [PyPI](https://pypi.org/project/viewmd/).

The active version is now written in Go — see the [root README](../README.md) for install instructions.

## Install (from PyPI)

```bash
pipx install viewmd
```

## Usage

```bash
viewmd          # Starts on port 8000
viewmd 3000     # Custom port
```

## Development

```bash
pip install -e ".[dev]"
pytest
pylint viewmd/ tests/
```
