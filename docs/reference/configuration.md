# Configuration

viewmd can be configured with a `.viewmd.yaml` file in the served directory.

## Config File

Create a `.viewmd.yaml` file in the root of the directory you're serving:

```yaml
show_all_files: true
```

## Options

| Key | Type | Description | Default |
|-----|------|-------------|---------|
| `show_all_files` | `boolean` | Show all files, not just markdown | `false` |

## Precedence

CLI flags override config file values. For example, if `.viewmd.yaml` sets `show_all_files: true` but you don't pass `-a`, all files will still be shown. The config file provides a default that applies when no flag is given.
