// Package markdown converts markdown text to HTML and serves markdown files.
package markdown

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/driangle/viewmd/apps/cli/internal/frontmatter"
	"github.com/driangle/viewmd/apps/cli/internal/render"
)

var converter = goldmark.New(
	goldmark.WithExtensions(
		extension.Table,
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
	),
)

// Convert transforms markdown text into an HTML string.
func Convert(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := converter.Convert([]byte(markdown), &buf); err != nil {
		return "", fmt.Errorf("goldmark conversion: %w", err)
	}
	return buf.String(), nil
}

// ServeMarkdown reads a markdown file, parses frontmatter, converts the body
// to HTML, and writes a full rendered page to w.
func ServeMarkdown(w http.ResponseWriter, filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
		return
	}

	meta, body := frontmatter.Parse(string(content))

	bodyHTML, err := Convert(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering markdown: %v", err), http.StatusInternalServerError)
		return
	}

	baseURL := computeBaseURL(filePath)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := render.RenderMarkdownPage(w, filepath.Base(filePath), meta, bodyHTML, baseURL); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering page: %v", err), http.StatusInternalServerError)
	}
}

// computeBaseURL returns "/" when the file is in the current directory,
// or "/dir/" for files in subdirectories.
func computeBaseURL(filePath string) string {
	dir := filepath.Dir(filePath)
	if dir == "." {
		return "/"
	}
	return "/" + filepath.ToSlash(dir) + "/"
}
