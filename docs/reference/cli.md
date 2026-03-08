# CLI Flags

## Usage

```
viewmd [flags] [port]
```

## Arguments

| Argument | Description | Default |
|----------|-------------|---------|
| `port` | Port number to serve on | `8000` |

## Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--show-all` | `-a` | Show all files, not just markdown | `false` |
| `--auto-readme` | | Auto-render README.md in directories | `false` |
| `--version` | | Print version and exit | |

## Examples

```bash
# Serve current directory on default port
viewmd

# Serve on port 3000
viewmd 3000

# Show all files
viewmd -a

# Show all files with auto-README on port 9000
viewmd -a --auto-readme 9000
```

## Exit

Press `Ctrl+C` to stop the server. viewmd handles `SIGINT` and `SIGTERM` for graceful shutdown.
