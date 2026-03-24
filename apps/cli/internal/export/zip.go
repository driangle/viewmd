// Package export provides file export functionality such as ZIP archive generation.
package export

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/driangle/viewmd/apps/cli/internal/frontmatter"
	"github.com/driangle/viewmd/apps/cli/internal/markdown"
	"github.com/driangle/viewmd/apps/cli/internal/render"
)

// IsIgnoredFunc checks whether a file (by basename and relative path) should be excluded.
type IsIgnoredFunc func(name, relPath string) bool

// WriteZip writes a ZIP archive to w containing all files under dirPath recursively.
// Markdown files (.md, .markdown) are converted to HTML; other files are included as-is.
// Each directory gets an index.html with a navigable file listing.
// When markdownOnly is true, only markdown files are included (matching viewmd's default mode).
func WriteZip(w io.Writer, dirPath string, isIgnored IsIgnoredFunc, markdownOnly bool, theme string) error {
	zw := zip.NewWriter(w)
	defer zw.Close()

	rootName := filepath.Base(dirPath)
	if rootName == "." || rootName == "/" {
		rootName = "files"
	}

	return writeDir(zw, dirPath, ".", rootName, isIgnored, markdownOnly, theme)
}

// writeDir recursively processes a directory, adding files and an index.html to the ZIP.
func writeDir(zw *zip.Writer, absDir, relDir, displayName string, isIgnored IsIgnoredFunc, markdownOnly bool, theme string) error {
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return err
	}

	// Build set of existing non-markdown filenames to detect collisions
	existingNames := make(map[string]bool)
	for _, e := range entries {
		if !e.IsDir() {
			existingNames[e.Name()] = true
		}
	}

	// Collect visible entries
	var items []render.DirEntry
	var dirs []os.DirEntry

	for _, e := range entries {
		name := e.Name()
		entryRel := name
		if relDir != "." {
			entryRel = relDir + "/" + name
		}

		if isIgnored != nil && isIgnored(name, entryRel) {
			continue
		}

		if e.IsDir() {
			if markdownOnly && !hasMarkdownFiles(filepath.Join(absDir, name)) {
				continue
			}
			items = append(items, render.DirEntry{
				Name:  name,
				Href:  name + "/_directory.html",
				IsDir: true,
			})
			dirs = append(dirs, e)
		} else {
			if markdownOnly && !isMarkdownFile(name) {
				continue
			}
			href := name
			if isMarkdownFile(name) {
				htmlName := replaceExt(name, ".html")
				if !existingNames[htmlName] {
					href = htmlName
				}
				// If collision, href stays as the original .md name
			}
			items = append(items, render.DirEntry{
				Name:  name,
				Href:  href,
				IsDir: false,
			})
		}
	}

	// Sort: directories first, then files, alphabetically
	sort.Slice(items, func(i, j int) bool {
		if items[i].IsDir != items[j].IsDir {
			return items[i].IsDir
		}
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})

	// Generate index.html for this directory
	if err := addIndexPage(zw, relDir, displayName, items, theme); err != nil {
		return err
	}

	// Add files
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		entryRel := name
		if relDir != "." {
			entryRel = relDir + "/" + name
		}

		if isIgnored != nil && isIgnored(name, entryRel) {
			continue
		}
		if markdownOnly && !isMarkdownFile(name) {
			continue
		}

		absPath := filepath.Join(absDir, name)
		canConvert := isMarkdownFile(name) && !existingNames[replaceExt(name, ".html")]
		if canConvert {
			if err := addMarkdownFile(zw, absPath, entryRel, theme); err != nil {
				return err
			}
		} else {
			if err := addRawFile(zw, absPath, entryRel); err != nil {
				return err
			}
		}
	}

	// Recurse into subdirectories
	for _, d := range dirs {
		childAbs := filepath.Join(absDir, d.Name())
		childRel := d.Name()
		if relDir != "." {
			childRel = relDir + "/" + d.Name()
		}
		if err := writeDir(zw, childAbs, childRel, d.Name(), isIgnored, markdownOnly, theme); err != nil {
			return err
		}
	}

	return nil
}

// addIndexPage renders a directory listing and writes it as index.html in the ZIP.
func addIndexPage(zw *zip.Writer, relDir, displayName string, items []render.DirEntry, theme string) error {
	breadcrumbs := buildExportBreadcrumbs(relDir, displayName)

	parentHref := ""
	if relDir != "." {
		parentHref = "../_directory.html"
	}

	var buf bytes.Buffer
	if err := render.RenderExportDirectoryPage(&buf, displayName, parentHref, items, breadcrumbs, theme); err != nil {
		return err
	}

	indexPath := "_directory.html"
	if relDir != "." {
		indexPath = relDir + "/_directory.html"
	}

	fw, err := zw.Create(indexPath)
	if err != nil {
		return err
	}
	_, err = fw.Write(buf.Bytes())
	return err
}

// buildExportBreadcrumbs creates breadcrumb segments for an exported directory.
// Each parent links to its relative ../index.html.
func buildExportBreadcrumbs(relDir, displayName string) []render.BreadcrumbSegment {
	if relDir == "." {
		return []render.BreadcrumbSegment{{Name: displayName}}
	}

	parts := strings.Split(relDir, "/")
	var crumbs []render.BreadcrumbSegment

	// Each ancestor gets a relative link going up the right number of levels
	for i, part := range parts {
		if i == len(parts)-1 {
			// Current directory — no link
			crumbs = append(crumbs, render.BreadcrumbSegment{Name: part})
		} else {
			// Link goes up (len(parts)-1-i) levels
			ups := strings.Repeat("../", len(parts)-1-i)
			crumbs = append(crumbs, render.BreadcrumbSegment{
				Name: part,
				Href: ups + "_directory.html",
			})
		}
	}
	return crumbs
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
	if err := render.RenderExportPage(&buf, htmlName, meta, bodyHTML, theme, "_directory.html"); err != nil {
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

// hasMarkdownFiles returns true if the directory contains at least one
// markdown file, searching recursively.
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
