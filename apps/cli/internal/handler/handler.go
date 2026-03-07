// Package handler provides the HTTP request handler for serving
// markdown files, text files, and directory listings.
package handler

import (
	"fmt"
	"html"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/driangle/viewmd/apps/cli/internal/classify"
	"github.com/driangle/viewmd/apps/cli/internal/markdown"
	"github.com/driangle/viewmd/apps/cli/internal/render"
)

// Handler serves markdown files, text files, and directory listings
// from a root directory.
type Handler struct {
	root       string
	AutoReadme bool
}

// New creates a Handler that serves files from the given root directory.
func New(root string) *Handler {
	return &Handler{root: root}
}

// ServeHTTP dispatches requests to the appropriate rendering path.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parsed, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	reqPath, err := url.PathUnescape(parsed.Path)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	reqPath = strings.TrimPrefix(reqPath, "/")
	reqPath = strings.TrimSuffix(reqPath, "/")

	timestamp := time.Now().Format("15:04:05")
	displayPath := reqPath
	if displayPath == "" {
		displayPath = "/"
	}
	fmt.Printf("[%s] Request: %s\n", timestamp, displayPath)

	if reqPath == "" {
		h.serveDirectoryListing(w, r, ".")
		return
	}

	fullPath := filepath.Join(h.root, reqPath)
	info, err := os.Stat(fullPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if info.IsDir() {
		if h.AutoReadme {
			readme := filepath.Join(fullPath, "README.md")
			if _, err := os.Stat(readme); err == nil {
				baseURL := "/" + reqPath + "/"
				parentHref := parentHrefFromPath(reqPath)
				markdown.ServeMarkdown(w, readme, baseURL, parentHref)
				return
			}
		}
		h.serveDirectoryListing(w, r, reqPath)
		return
	}

	parentHref := parentHrefFromPath(reqPath)
	ext := strings.ToLower(filepath.Ext(fullPath))
	switch {
	case ext == ".md" || ext == ".markdown":
		markdown.ServeMarkdown(w, fullPath, parentHref, parentHref)
	case classify.IsTextFile(info.Name()):
		serveTextFile(w, r, fullPath, parentHref)
	default:
		http.ServeFile(w, r, fullPath)
	}
}

// serveTextFile reads a file as UTF-8 text and renders it as an HTML page.
// parentHrefFromPath returns the URL path to the parent directory listing.
// For "foo/bar/baz.md" it returns "/foo/bar/", for "baz.md" it returns "/".
func parentHrefFromPath(reqPath string) string {
	dir := filepath.Dir(reqPath)
	if dir == "." {
		return "/"
	}
	return "/" + filepath.ToSlash(dir) + "/"
}

func serveTextFile(w http.ResponseWriter, r *http.Request, filePath string, parentHref string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err),
			http.StatusInternalServerError)
		return
	}

	if !utf8.Valid(content) {
		http.ServeFile(w, r, filePath)
		return
	}

	escaped := html.EscapeString(string(content))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := render.RenderTextPage(w, filepath.Base(filePath), escaped, parentHref); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering page: %v", err),
			http.StatusInternalServerError)
	}
}
