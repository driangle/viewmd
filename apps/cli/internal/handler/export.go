package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/driangle/viewmd/apps/cli/internal/export"
	"github.com/driangle/viewmd/apps/cli/internal/frontmatter"
	"github.com/driangle/viewmd/apps/cli/internal/markdown"
	"github.com/driangle/viewmd/apps/cli/internal/render"
)

// serveZipExport generates a ZIP archive of the directory at fullPath and
// streams it as a download. Markdown files are converted to HTML; all other
// files are included as-is.
func (h *Handler) serveZipExport(w http.ResponseWriter, r *http.Request, fullPath string, reqPath string) {
	zipName := zipFileName(reqPath, h.root)

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, zipName))

	patterns := h.IgnorePatterns
	isIgnored := func(name, relPath string) bool {
		if len(patterns) == 0 {
			return false
		}
		return isIgnored(name, relPath, patterns)
	}

	theme := r.URL.Query().Get("theme")
	if err := export.WriteZip(w, fullPath, isIgnored, !h.ShowAll, theme); err != nil {
		http.Error(w, fmt.Sprintf("Error creating ZIP: %v", err), http.StatusInternalServerError)
	}
}

// serveHTMLExport renders a single markdown file as a self-contained HTML page
// and streams it as a download.
func serveHTMLExport(w http.ResponseWriter, r *http.Request, fullPath string) {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
		return
	}

	meta, body := frontmatter.Parse(string(content))

	bodyHTML, err := markdown.Convert(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering markdown: %v", err), http.StatusInternalServerError)
		return
	}

	fileName := filepath.Base(fullPath)
	htmlName := replaceExtToHTML(fileName)

	theme := r.URL.Query().Get("theme")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, htmlName))

	if err := render.RenderExportPage(w, htmlName, meta, bodyHTML, theme, ""); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering page: %v", err), http.StatusInternalServerError)
	}
}

// replaceExtToHTML replaces the file extension with .html.
func replaceExtToHTML(name string) string {
	ext := filepath.Ext(name)
	return strings.TrimSuffix(name, ext) + ".html"
}

// zipFileName derives the ZIP archive filename from the request path.
// For root directories it uses the root directory's basename.
func zipFileName(reqPath string, root string) string {
	if reqPath == "" || reqPath == "." {
		name := filepath.Base(root)
		if name == "." || name == "/" {
			name = "archive"
		}
		return name + ".zip"
	}
	return filepath.Base(reqPath) + ".zip"
}
