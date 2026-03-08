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
| `--watch` | `-w` | Watch for file changes and auto-reload browser | `false` |
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

# Watch for file changes and auto-reload the browser
viewmd --watch

# Combine watch with other flags
viewmd -a -w --auto-readme
```

## Exit

Press `Ctrl+C` to stop the server. viewmd handles `SIGINT` and `SIGTERM` for graceful shutdown.
