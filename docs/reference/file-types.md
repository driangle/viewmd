# Supported File Types

viewmd classifies files into categories for rendering. Files not recognized as text are treated as binary and shown with a download link.

## Markdown

Rendered as full HTML with styled typography.

| Extension |
|-----------|
| `.md` |
| `.markdown` |

## Text / Code

Displayed with syntax highlighting in a monospace font.

### By Extension

`.txt`, `.log`, `.json`, `.xml`, `.yaml`, `.yml`, `.toml`, `.ini`, `.cfg`, `.conf`, `.sh`, `.bash`, `.zsh`, `.fish`, `.py`, `.js`, `.ts`, `.jsx`, `.tsx`, `.java`, `.c`, `.cpp`, `.h`, `.hpp`, `.cs`, `.go`, `.rs`, `.rb`, `.php`, `.swift`, `.kt`, `.sql`, `.html`, `.css`, `.scss`, `.sass`, `.less`, `.vue`, `.svelte`, `.r`, `.m`, `.scala`, `.pl`, `.lua`, `.vim`, `.el`, `.clj`, `.ex`, `.exs`, `.dockerfile`, `.env`, `.gitignore`, `.gitattributes`, `.editorconfig`, `.eslintrc`, `.prettierrc`, `.babelrc`

### By Filename

These filenames are recognized regardless of extension (case-insensitive):

`Makefile`, `Dockerfile`, `Gemfile`, `Rakefile`, `Procfile`, `Jenkinsfile`, `LICENSE`, `README`, `CHANGELOG`, `AUTHORS`, `CONTRIBUTORS`, `CODEOWNERS`

## Auto-Detection

Files with unrecognized extensions are probed for UTF-8 validity (first 8 KB). Valid UTF-8 files are rendered as plain text; invalid files are treated as binary.

## Binary / Unsupported

Files that can't be previewed show a "No preview available" page with a download link. Append `?raw=1` to any file URL to download it directly.
