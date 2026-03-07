# viewmd: Python to Go Migration Plan

## Goal

Rewrite viewmd as a Go CLI that achieves feature parity with the current Python implementation, then serves as the foundation for future enhancements.

## Current Feature Set (parity target)

1. **CLI** — `viewmd [port]`, defaults to port 8000
2. **Markdown rendering** — `.md`/`.markdown` files converted to styled HTML (fenced code, tables, nl2br)
3. **Frontmatter** — YAML-like `---` delimited key:value pairs parsed and rendered as an HTML table above the content
4. **Text file viewer** — Known extensions/filenames/dotfiles displayed in `<pre>` with HTML escaping; falls back to binary on decode errors
5. **Directory listing** — Sorted (dirs first), auto-redirects to `README.md` if present, parent link for subdirs
6. **Binary passthrough** — Unknown file types served as-is (images, PDFs, etc.)
7. **Graceful errors** — Port-in-use detection, Ctrl+C shutdown, 404 for missing files

## Architecture

```
cmd/viewmd/main.go      — Entry point, CLI arg parsing, server startup
internal/handler/        — HTTP handler (routing, response helpers)
internal/frontmatter/    — parse_frontmatter equivalent
internal/render/         — HTML template rendering (markdown page, text page, directory page)
internal/classify/       — File type classification (text extensions, known filenames, dotfiles)
```

### Key Dependencies

| Purpose              | Go Package                                      |
|----------------------|-------------------------------------------------|
| Markdown → HTML      | `github.com/yuin/goldmark` (fenced code, tables) |
| HTML templates       | `html/template` (stdlib)                         |
| HTTP server          | `net/http` (stdlib)                              |
| CLI                  | `flag` or bare `os.Args` (stdlib)                |

### Why goldmark

goldmark is the de facto Go markdown library. It supports fenced code blocks and tables via built-in extensions, matching the Python `markdown` library features we use.

## Migration Phases

### Phase 1: Project Scaffolding

- [ ] Initialize Go module (`go mod init github.com/driangle/viewmd`)
- [ ] Set up directory structure (`cmd/`, `internal/`)
- [ ] Add `Makefile` with `build`, `test`, `lint` targets
- [ ] Add goldmark dependency

### Phase 2: Frontmatter Parser

- [ ] Implement `frontmatter.Parse(content string) (map[string]string, string)` matching current Python behavior:
  - Must start with `---\n`
  - Split on `---` taking first two delimiters only
  - Lines with `:` become key:value (trimmed), lines without are skipped
  - Returns nil map (not empty) when no frontmatter found
- [ ] Port all `TestParseFrontmatter` cases as Go table-driven tests

### Phase 3: File Classification

- [ ] Implement `classify.IsTextFile(name string) bool` with the same extension/filename/dotfile sets
- [ ] Port `TestIsTextFile` cases

### Phase 4: HTML Rendering

- [ ] Create Go `html/template` templates for:
  - Markdown page (with frontmatter table, base href, CSS)
  - Text file page (filename header, pre block, CSS)
  - Directory listing page (sorted entries, parent link, CSS)
- [ ] Ensure CSS output matches current Python templates exactly
- [ ] Templates embedded via `embed.FS`

### Phase 5: HTTP Handler

- [ ] Implement request routing in `handler.Handler`:
  - URL decode path
  - Route: root → directory listing
  - Route: directory with README.md → serve markdown
  - Route: directory without README.md → directory listing
  - Route: `.md`/`.markdown` → render markdown
  - Route: text file → render text (with binary fallback on invalid UTF-8)
  - Route: other files → `http.FileServer` passthrough
  - 404 for missing files
- [ ] Request logging with timestamp `[HH:MM:SS] Request: /path`

### Phase 6: CLI & Server Entry Point

- [ ] Parse optional port argument (default 8000)
- [ ] Bind server, detect port-in-use, print banner
- [ ] Graceful shutdown on SIGINT/SIGTERM
- [ ] Version string

### Phase 7: Integration Tests

- [ ] Port `TestServeMarkdownFrontmatter` tests (real HTTP server on random port)
- [ ] Port `TestServeTextFile` tests
- [ ] Port `TestServeDirectoryListing` tests
- [ ] Port `TestRouting` tests
- [ ] Use `testing` + `net/http/httptest` or real server on `:0`
- [ ] Copy fixture `.md` files into temp dirs, same pattern as Python tests

### Phase 8: Build & Release

- [ ] Update GitHub Actions for Go build/test
- [ ] Cross-compile binaries (linux/darwin/windows, amd64/arm64)
- [ ] Create GitHub Release with binaries (replaces PyPI publishing)
- [ ] Update README with new install instructions (`go install` or binary download)
- [ ] Update CLAUDE.md project description

### Phase 9: Cleanup

- [ ] Remove Python source (`viewmd/`, `tests/`, `pyproject.toml`, `requirements.txt`, etc.)
- [ ] Remove Python CI config
- [ ] Keep `LICENSE`, `README.md`

## Design Decisions

- **`html/template` over string concatenation** — safer (auto-escaping), cleaner separation of markup from logic.
- **`embed.FS` for templates** — single binary, no runtime file dependencies.
- **Stdlib HTTP** — no need for a framework; `net/http` is more than sufficient for a file server.
- **goldmark** — well-maintained, extensible, supports the extensions we need out of the box.
- **No nl2br** — goldmark doesn't have a direct nl2br equivalent. We'll evaluate whether this is needed (it's unusual for markdown renderers) and implement a simple AST transformer if required.

## Out of Scope (future phases)

- Live reload / file watching
- Syntax highlighting in code blocks
- Custom themes / CSS
- Search
- Table of contents sidebar
