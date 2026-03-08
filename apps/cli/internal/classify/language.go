package classify

import (
	"path/filepath"
	"strings"
)

// langByExtension maps file extensions to highlight.js language identifiers.
var langByExtension = map[string]string{
	".go":         "go",
	".py":         "python",
	".js":         "javascript",
	".ts":         "typescript",
	".jsx":        "javascript",
	".tsx":        "typescript",
	".java":       "java",
	".c":          "c",
	".cpp":        "cpp",
	".h":          "c",
	".hpp":        "cpp",
	".cs":         "csharp",
	".rs":         "rust",
	".rb":         "ruby",
	".php":        "php",
	".swift":      "swift",
	".kt":         "kotlin",
	".scala":      "scala",
	".pl":         "perl",
	".lua":        "lua",
	".r":          "r",
	".m":          "objectivec",
	".clj":        "clojure",
	".ex":         "elixir",
	".exs":        "elixir",
	".el":         "lisp",
	".vim":        "vim",
	".sh":         "bash",
	".bash":       "bash",
	".zsh":        "bash",
	".fish":       "shell",
	".sql":        "sql",
	".html":       "xml",
	".xml":        "xml",
	".css":        "css",
	".scss":       "scss",
	".sass":       "scss",
	".less":       "less",
	".json":       "json",
	".yaml":       "yaml",
	".yml":        "yaml",
	".toml":       "ini",
	".ini":        "ini",
	".cfg":        "ini",
	".conf":       "ini",
	".dockerfile": "dockerfile",
	".vue":        "xml",
	".svelte":     "xml",
	".txt":        "",
	".log":        "",
	".env":        "bash",
}

// langByFilename maps exact filenames (lowercase) to highlight.js language identifiers.
var langByFilename = map[string]string{
	"makefile":    "makefile",
	"dockerfile":  "dockerfile",
	"gemfile":     "ruby",
	"rakefile":    "ruby",
	"jenkinsfile": "groovy",
}

// DetectLanguage returns the highlight.js language identifier for a filename.
// Returns empty string if no language mapping is found.
func DetectLanguage(name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	if lang, ok := langByExtension[ext]; ok {
		return lang
	}
	base := strings.ToLower(filepath.Base(name))
	if lang, ok := langByFilename[base]; ok {
		return lang
	}
	return ""
}
