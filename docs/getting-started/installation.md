# Installation

## Homebrew (macOS / Linux)

```bash
brew install driangle/tap/viewmd
```

## From Source

Requires [Go](https://go.dev/) 1.22 or later.

```bash
go install github.com/driangle/viewmd/apps/cli/cmd/viewmd@latest
```

## Build Locally

```bash
git clone https://github.com/driangle/viewmd.git
cd viewmd/apps/cli
make build
# Binary at apps/cli/bin/viewmd
```

## Verify Installation

```bash
viewmd --version
```

## Python Version (Legacy)

The original Python version is available on [PyPI](https://pypi.org/project/viewmd/):

```bash
pipx install viewmd
```

The Python version is archived and no longer actively developed. The Go CLI is the recommended version.
