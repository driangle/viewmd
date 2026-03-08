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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := render.RenderDirectoryPage(w, dirPath, parentHref, items); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering page: %v", err),
			http.StatusInternalServerError)
	}
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
