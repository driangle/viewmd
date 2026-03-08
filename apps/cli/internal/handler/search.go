package handler

import (
	"bufio"
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/driangle/viewmd/apps/cli/internal/classify"
)

const maxSearchResults = 50
const maxFileSize = 1 << 20 // 1MB

// searchResult represents a single search match.
type searchResult struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	IsDir   bool   `json:"isDir"`
	Match   string `json:"match"`
	Snippet string `json:"snippet,omitempty"`
}

// serveSearch handles GET /-/search?q=...&mode=name|content|both
func (h *Handler) serveSearch(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "both"
	}

	w.Header().Set("Content-Type", "application/json")

	if q == "" {
		json.NewEncoder(w).Encode(map[string][]searchResult{"results": {}})
		return
	}

	qLower := strings.ToLower(q)
	var results []searchResult

	filepath.WalkDir(h.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || len(results) >= maxSearchResults {
			return nil
		}

		relPath, _ := filepath.Rel(h.root, path)
		if relPath == "." {
			return nil
		}

		// Skip hidden directories
		name := d.Name()
		if strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		// In non-ShowAll mode, only search markdown files and dirs with markdown content
		if !h.ShowAll {
			if d.IsDir() {
				if !hasMarkdownFiles(path) {
					return fs.SkipDir
				}
			} else if !isMarkdownFile(name) {
				return nil
			}
		}

		nameMatch := (mode == "name" || mode == "both") &&
			strings.Contains(strings.ToLower(name), qLower)

		var snippet string
		var contentMatch bool
		if (mode == "content" || mode == "both") && !d.IsDir() {
			snippet, contentMatch = searchFileContent(path, name, qLower)
		}

		if nameMatch || contentMatch {
			match := "name"
			if contentMatch && !nameMatch {
				match = "content"
			} else if contentMatch && nameMatch {
				match = "both"
			}
			results = append(results, searchResult{
				Path:    filepath.ToSlash(relPath),
				Name:    name,
				IsDir:   d.IsDir(),
				Match:   match,
				Snippet: snippet,
			})
		}

		return nil
	})

	json.NewEncoder(w).Encode(map[string][]searchResult{"results": results})
}

// searchFileContent scans a text file for a case-insensitive match
// and returns the matching line as a snippet.
func searchFileContent(path, name, qLower string) (string, bool) {
	if !isSearchableFile(name) {
		return "", false
	}

	info, err := os.Stat(path)
	if err != nil || info.Size() > maxFileSize {
		return "", false
	}

	f, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), qLower) {
			snippet := strings.TrimSpace(line)
			if len(snippet) > 120 {
				snippet = snippet[:120] + "..."
			}
			return snippet, true
		}
	}
	return "", false
}

// isSearchableFile returns true if the file is a text or markdown file.
func isSearchableFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	if ext == ".md" || ext == ".markdown" {
		return true
	}
	return classify.IsTextFile(name)
}
