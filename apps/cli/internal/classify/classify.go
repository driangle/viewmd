// Package classify determines file types for routing requests
// to the appropriate rendering path.
package classify

import (
	"path/filepath"
	"strings"
)

// textExtensions contains file extensions recognized as text files.
// Markdown extensions (.md, .markdown) are intentionally excluded
// because they have their own rendering path.
var textExtensions = map[string]bool{
	".txt": true, ".log": true, ".json": true, ".xml": true,
	".yaml": true, ".yml": true, ".toml": true, ".ini": true,
	".cfg": true, ".conf": true, ".sh": true, ".bash": true,
	".zsh": true, ".fish": true, ".py": true, ".js": true,
	".ts": true, ".jsx": true, ".tsx": true, ".java": true,
	".c": true, ".cpp": true, ".h": true, ".hpp": true,
	".cs": true, ".go": true, ".rs": true, ".rb": true,
	".php": true, ".swift": true, ".kt": true, ".sql": true,
	".html": true, ".css": true, ".scss": true, ".sass": true,
	".less": true, ".vue": true, ".svelte": true, ".r": true,
	".m": true, ".scala": true, ".pl": true, ".lua": true,
	".vim": true, ".el": true, ".clj": true, ".ex": true,
	".exs": true, ".dockerfile": true, ".env": true,
	".gitignore": true, ".gitattributes": true, ".editorconfig": true,
	".eslintrc": true, ".prettierrc": true, ".babelrc": true,
}

// textFilenames contains exact filenames (case-insensitive) recognized as text.
var textFilenames = map[string]bool{
	"makefile":     true,
	"dockerfile":   true,
	"gemfile":      true,
	"rakefile":     true,
	"procfile":     true,
	"jenkinsfile":  true,
	"license":      true,
	"readme":       true,
	"changelog":    true,
	"authors":      true,
	"contributors": true,
	"codeowners":   true,
}

// IsTextFile reports whether the given filename should be treated as a text
// file for rendering purposes. Markdown files (.md, .markdown) return false
// because they use a separate rendering path.
func IsTextFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	if ext == ".md" || ext == ".markdown" {
		return false
	}
	if textExtensions[ext] {
		return true
	}
	base := strings.ToLower(filepath.Base(name))
	if textFilenames[base] {
		return true
	}
	return strings.HasPrefix(filepath.Base(name), ".")
}
