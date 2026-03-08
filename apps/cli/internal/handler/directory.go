package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/driangle/viewmd/apps/cli/internal/render"
)

// isIgnored returns true if the entry matches any of the given glob patterns.
// Patterns without "/" match against the basename only.
// Patterns with "/" match against the full relative path.
// Patterns with "**" match across multiple path segments.
func isIgnored(name, relPath string, patterns []string) bool {
	for _, p := range patterns {
		if strings.Contains(p, "/") {
			if matchPathPattern(relPath, p) {
				return true
			}
		} else {
			if matched, _ := filepath.Match(p, name); matched {
				return true
			}
		}
	}
	return false
}

// matchPathPattern matches a relative path against a pattern that may contain "/".
// Supports "**" to match zero or more path segments.
func matchPathPattern(relPath, pattern string) bool {
	if !strings.Contains(pattern, "**") {
		// Exact path match with filepath.Match (supports single * within segments)
		if matched, _ := filepath.Match(pattern, relPath); matched {
			return true
		}
		// Also match as a prefix (e.g. pattern ".claude/worktrees" matches ".claude/worktrees/foo")
		if strings.HasPrefix(relPath, pattern+"/") {
			return true
		}
		return false
	}
	return matchDoublestar(relPath, pattern)
}

// matchDoublestar matches a path against a pattern containing "**".
// "**" matches zero or more path segments.
func matchDoublestar(path, pattern string) bool {
	parts := strings.SplitN(pattern, "**", 2)
	before, after := parts[0], parts[1]

	// Before "**" must match as a prefix
	if before != "" {
		if !strings.HasPrefix(path, before) {
			return false
		}
		path = path[len(before):]
	}

	// After "**" — if empty, everything matches
	if after == "" {
		return true
	}
	// Strip leading "/" from after since ** absorbs separators
	after = strings.TrimPrefix(after, "/")
	if after == "" {
		return true
	}

	// Try matching after against every suffix of the remaining path
	for i := 0; i <= len(path); i++ {
		if i > 0 && path[i-1] != '/' {
			continue
		}
		suffix := path[i:]
		if matched, _ := filepath.Match(after, suffix); matched {
			return true
		}
		// Also allow after to be a prefix (for nested matches)
		if strings.HasPrefix(suffix, after+"/") || suffix == after {
			return true
		}
	}
	return false
}

// serveDirectoryListing renders a sorted directory listing page.
// Directories are listed first (alphabetically), then files (alphabetically).
func (h *Handler) serveDirectoryListing(w http.ResponseWriter, _ *http.Request, dirPath string) {
	fullPath := filepath.Join(h.root, dirPath)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing directory: %v", err),
			http.StatusInternalServerError)
		return
	}

	totalCount := len(entries)

	if len(h.IgnorePatterns) > 0 {
		entries = filterIgnoredEntries(entries, dirPath, h.IgnorePatterns)
	}

	if !h.ShowAll {
		entries = filterMarkdownEntries(fullPath, entries)
	}

	sort.Slice(entries, func(i, j int) bool {
		iDir := entries[i].IsDir()
		jDir := entries[j].IsDir()
		if iDir != jDir {
			return iDir
		}
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	items := make([]render.DirEntry, 0, len(entries))
	for _, e := range entries {
		href := e.Name()
		if dirPath != "." {
			href = dirPath + "/" + e.Name()
		}
		items = append(items, render.DirEntry{
			Name:  e.Name(),
			Href:  href,
			IsDir: e.IsDir(),
		})
	}

	var parentHref *string
	if dirPath != "." {
		parent := filepath.Dir(dirPath)
		if parent == "." {
			parent = ""
		}
		parentHref = &parent
	}

	bcPath := dirPath
	if bcPath == "." {
		bcPath = ""
	}
	breadcrumbs := render.BuildBreadcrumbs(bcPath, h.root)

	var emptyReason string
	if len(items) == 0 {
		if totalCount == 0 {
			emptyReason = "empty"
		} else {
			emptyReason = "all_hidden"
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := render.RenderDirectoryPage(w, dirPath, parentHref, items, breadcrumbs, emptyReason); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering page: %v", err),
			http.StatusInternalServerError)
	}
}

// filterIgnoredEntries removes entries whose names match any ignore pattern.
func filterIgnoredEntries(entries []os.DirEntry, dirPath string, patterns []string) []os.DirEntry {
	filtered := make([]os.DirEntry, 0, len(entries))
	for _, e := range entries {
		relPath := e.Name()
		if dirPath != "." {
			relPath = dirPath + "/" + e.Name()
		}
		if !isIgnored(e.Name(), relPath, patterns) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// filterMarkdownEntries returns only markdown files and directories
// that contain at least one markdown file recursively.
func filterMarkdownEntries(dirPath string, entries []os.DirEntry) []os.DirEntry {
	filtered := make([]os.DirEntry, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			if hasMarkdownFiles(filepath.Join(dirPath, e.Name())) {
				filtered = append(filtered, e)
			}
		} else if isMarkdownFile(e.Name()) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// isMarkdownFile returns true if the filename has a .md or .markdown extension.
func isMarkdownFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".md" || ext == ".markdown"
}

// hasMarkdownFiles returns true if the directory at path contains
// at least one markdown file, searching recursively.
func hasMarkdownFiles(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() {
			if hasMarkdownFiles(filepath.Join(path, e.Name())) {
				return true
			}
		} else if isMarkdownFile(e.Name()) {
			return true
		}
	}
	return false
}
