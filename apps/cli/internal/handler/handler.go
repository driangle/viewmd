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
	"unicode/utf8"

	"github.com/driangle/viewmd/apps/cli/internal/classify"
	"github.com/driangle/viewmd/apps/cli/internal/markdown"
	"github.com/driangle/viewmd/apps/cli/internal/render"
)

// Handler serves markdown files, text files, and directory listings
// from a root directory.
type Handler struct {
	root           string
	AutoReadme     bool
	ShowAll        bool
	IgnorePatterns []string
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

	if reqPath == "-/search" {
		h.serveSearch(w, r)
		return
	}

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
				breadcrumbs := render.BuildBreadcrumbs(reqPath+"/README.md", h.root)
				markdown.ServeMarkdown(w, readme, baseURL, parentHref, breadcrumbs)
				return
			}
		}
		h.serveDirectoryListing(w, r, reqPath)
		return
	}

	if r.URL.Query().Get("raw") == "1" {
		http.ServeFile(w, r, fullPath)
		return
	}

	parentHref := parentHrefFromPath(reqPath)
	breadcrumbs := render.BuildBreadcrumbs(reqPath, h.root)
	ext := strings.ToLower(filepath.Ext(fullPath))
	switch {
	case ext == ".md" || ext == ".markdown":
		markdown.ServeMarkdown(w, fullPath, parentHref, parentHref, breadcrumbs)
	case classify.IsImageFile(info.Name()):
		serveImageFile(w, info, parentHref, breadcrumbs)
	case classify.IsTextFile(info.Name()):
		serveTextFile(w, fullPath, parentHref, breadcrumbs)
	default:
		serveUnknownFile(w, fullPath, info, parentHref, breadcrumbs)
	}
}

// parentHrefFromPath returns the URL path to the parent directory listing.
// For "foo/bar/baz.md" it returns "/foo/bar/", for "baz.md" it returns "/".
func parentHrefFromPath(reqPath string) string {
	dir := filepath.Dir(reqPath)
	if dir == "." {
		return "/"
	}
	return "/" + filepath.ToSlash(dir) + "/"
}

// serveImageFile renders the image viewer page.
func serveImageFile(w http.ResponseWriter, info os.FileInfo, parentHref string, breadcrumbs []render.BreadcrumbSegment) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := render.RenderImagePage(w, info.Name(), info.Size(), parentHref, breadcrumbs); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering page: %v", err), http.StatusInternalServerError)
	}
}

// serveUnknownFile handles files with unrecognized extensions by probing the
// first 8KB for valid UTF-8. Text files are rendered with the text template;
// binary files get an unsupported-file page with a download link.
func serveUnknownFile(w http.ResponseWriter, filePath string, info os.FileInfo, parentHref string, breadcrumbs []render.BreadcrumbSegment) {
	f, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	buf := make([]byte, 8192)
	n, _ := f.Read(buf)
	if n > 0 && utf8.Valid(buf[:n]) {
		serveTextFile(w, filePath, parentHref, breadcrumbs)
		return
	}
	serveUnsupportedFile(w, info, parentHref, breadcrumbs)
}

// serveUnsupportedFile renders the unsupported file page with a download link.
func serveUnsupportedFile(w http.ResponseWriter, info os.FileInfo, parentHref string, breadcrumbs []render.BreadcrumbSegment) {
	ext := filepath.Ext(info.Name())
	fileType := strings.TrimPrefix(ext, ".")
	downloadHref := "?raw=1"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := render.RenderUnsupportedPage(w, info.Name(), fileType, info.Size(), downloadHref, parentHref, breadcrumbs); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering page: %v", err), http.StatusInternalServerError)
	}
}

// serveTextFile reads a file as UTF-8 text and renders it as an HTML page.
func serveTextFile(w http.ResponseWriter, filePath string, parentHref string, breadcrumbs []render.BreadcrumbSegment) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err),
			http.StatusInternalServerError)
		return
	}

	if !utf8.Valid(content) {
		info, err := os.Stat(filePath)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		serveUnsupportedFile(w, info, parentHref, breadcrumbs)
		return
	}

	raw := string(content)
	escaped := html.EscapeString(raw)
	lang := classify.DetectLanguage(filepath.Base(filePath))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := render.RenderTextPage(w, filepath.Base(filePath), escaped, parentHref, raw, lang, breadcrumbs); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering page: %v", err),
			http.StatusInternalServerError)
	}
}
