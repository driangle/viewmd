# Search

viewmd includes a built-in full-text search feature accessible from any directory listing.

## Using Search

Press `/` to open the search panel, then type your query. Results appear as you type.

## Search Modes

Toggle between modes with `Tab`:

| Mode | Description |
|------|-------------|
| **Both** (default) | Search filenames and file contents |
| **Name** | Search filenames only |
| **Content** | Search file contents only |

## Search Behavior

- **Case-insensitive** matching
- Returns up to **50 results**
- Skips files larger than **1 MB**
- Skips hidden directories (starting with `.`)
- Shows a content snippet (first matching line, up to 120 characters) for content matches

## Filtering

Search respects the current display mode:

- **Default mode**: only searches markdown files
- **Show-all mode** (`-a`): searches all text and markdown files
