// Package export provides file export functionality such as ZIP archive generation.
package export

import (
	"archive/zip"
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/driangle/viewmd/apps/cli/internal/frontmatter"
	"github.com/driangle/viewmd/apps/cli/internal/markdown"
	"github.com/driangle/viewmd/apps/cli/internal/render"
)

// IsIgnoredFunc checks whether a file (by basename and relative path) should be excluded.
type IsIgnoredFunc func(name, relPath string) bool

// WriteZip writes a ZIP archive to w containing all files under dirPath recursively.
// Markdown files (.md, .markdown) are converted to HTML; other files are included as-is.
// When markdownOnly is true, only markdown files are included (matching viewmd's default mode).
// Relative paths in the archive are preserved from dirPath.
func WriteZip(w io.Writer, dirPath string, isIgnored IsIgnoredFunc, markdownOnly bool, theme string) error {
	zw := zip.NewWriter(w)
	defer zw.Close()

	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		if isIgnored != nil && isIgnored(d.Name(), relPath) {
			return nil
		}

		if markdownOnly && !isMarkdownFile(d.Name()) {
			return nil
		}

		if isMarkdownFile(d.Name()) {
			return addMarkdownFile(zw, path, relPath, theme)
		}
		return addRawFile(zw, path, relPath)
	})
}

// addMarkdownFile reads a markdown file, renders it as a self-contained HTML page
// with full viewmd styling, and writes it to the ZIP with an .html extension.
func addMarkdownFile(zw *zip.Writer, filePath, relPath, theme string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	meta, body := frontmatter.Parse(string(content))

	bodyHTML, err := markdown.Convert(body)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	htmlName := replaceExt(filepath.Base(filePath), ".html")
	if err := render.RenderExportPage(&buf, htmlName, meta, bodyHTML, theme); err != nil {
		return err
	}

	htmlPath := replaceExt(relPath, ".html")
	fw, err := zw.Create(htmlPath)
	if err != nil {
		return err
	}
	_, err = fw.Write(buf.Bytes())
	return err
}

// addRawFile copies a file as-is into the ZIP archive.
func addRawFile(zw *zip.Writer, filePath, relPath string) error {
	fw, err := zw.Create(relPath)
	if err != nil {
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(fw, f)
	return err
}

// isMarkdownFile returns true if the filename has a .md or .markdown extension.
func isMarkdownFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".md" || ext == ".markdown"
}

// replaceExt replaces the file extension in path with newExt.
func replaceExt(path, newExt string) string {
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)] + newExt
}
