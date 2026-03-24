// Package render provides HTML rendering for markdown content,
// text files, and directory listings.
package render

import (
	"embed"
	"fmt"
	"html/template"
	"io"

	"github.com/driangle/viewmd/apps/cli/internal/frontmatter"
)

// Version is set by the main package at startup.
var Version string

// WatchMode is set by the main package when --watch is enabled.
var WatchMode bool

//go:embed templates/*.html
var templateFS embed.FS

var templates = template.Must(
	template.ParseFS(templateFS, "templates/*.html"),
)

// RenderMarkdownPage writes a full HTML page for rendered markdown content.
// meta may be nil if no frontmatter was parsed.
// bodyHTML is already-rendered HTML from the markdown converter.
func RenderMarkdownPage(w io.Writer, fileName string, meta []frontmatter.KeyValue, bodyHTML string, baseURL string, parentHref string, rawContent string, breadcrumbs []BreadcrumbSegment) error {
	data := markdownData{
		FileName:        fileName,
		BaseURL:         baseURL,
		ParentHref:      parentHref,
		Breadcrumbs:     breadcrumbs,
		FrontmatterRows: meta,
		BodyHTML:        template.HTML(bodyHTML),
		RawContent:      rawContent,
		Version:         Version,
		WatchMode:       WatchMode,
	}
	return templates.ExecuteTemplate(w, "markdown.html", data)
}

// RenderTextPage writes a full HTML page for a plain text file.
// escapedContent should already be HTML-escaped.
// language is the highlight.js language identifier (empty for no highlighting).
func RenderTextPage(w io.Writer, fileName string, escapedContent string, parentHref string, rawContent string, language string, breadcrumbs []BreadcrumbSegment) error {
	data := textData{
		FileName:    fileName,
		ParentHref:  parentHref,
		Breadcrumbs: breadcrumbs,
		Content:     template.HTML(escapedContent),
		RawContent:  rawContent,
		Language:    language,
		Version:     Version,
		WatchMode:   WatchMode,
	}
	return templates.ExecuteTemplate(w, "text.html", data)
}

// RenderUnsupportedPage writes a full HTML page for an unsupported file type.
func RenderUnsupportedPage(w io.Writer, fileName string, fileType string, fileSize int64, downloadHref string, parentHref string, breadcrumbs []BreadcrumbSegment) error {
	data := unsupportedData{
		FileName:     fileName,
		FileType:     fileType,
		FileSize:     formatFileSize(fileSize),
		DownloadHref: downloadHref,
		ParentHref:   parentHref,
		Breadcrumbs:  breadcrumbs,
		Version:      Version,
		WatchMode:    WatchMode,
	}
	return templates.ExecuteTemplate(w, "unsupported.html", data)
}

// formatFileSize returns a human-readable file size string.
func formatFileSize(size int64) string {
	switch {
	case size < 1024:
		return fmt.Sprintf("%d B", size)
	case size < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
	default:
		return fmt.Sprintf("%.1f GB", float64(size)/(1024*1024*1024))
	}
}

// RenderImagePage writes a full HTML page for an inline image viewer.
func RenderImagePage(w io.Writer, fileName string, fileSize int64, parentHref string, breadcrumbs []BreadcrumbSegment) error {
	data := imageData{
		FileName:    fileName,
		ImageSrc:    "?raw=1",
		FileSize:    formatFileSize(fileSize),
		ParentHref:  parentHref,
		Breadcrumbs: breadcrumbs,
		Version:     Version,
		WatchMode:   WatchMode,
	}
	return templates.ExecuteTemplate(w, "image.html", data)
}

// RenderNotFoundPage writes a styled 404 page.
func RenderNotFoundPage(w io.Writer, path string, parentHref string, breadcrumbs []BreadcrumbSegment) error {
	data := notFoundData{
		Path:        path,
		ParentHref:  parentHref,
		Breadcrumbs: breadcrumbs,
		Version:     Version,
	}
	return templates.ExecuteTemplate(w, "notfound.html", data)
}

// RenderExportPage writes a self-contained HTML page for a markdown file,
// suitable for offline viewing in a browser. Includes all styling inline
// with no server dependencies. Theme should be "dark" or "light".
func RenderExportPage(w io.Writer, fileName string, meta []frontmatter.KeyValue, bodyHTML string, theme string) error {
	if theme != "dark" {
		theme = "light"
	}
	data := exportData{
		FileName:        fileName,
		FrontmatterRows: meta,
		BodyHTML:        template.HTML(bodyHTML),
		Theme:           theme,
	}
	return templates.ExecuteTemplate(w, "export.html", data)
}

// RenderDirectoryPage writes a full HTML page for a directory listing.
// parentHref is nil for the root directory.
// emptyReason is "" when items exist, "empty" for truly empty dirs,
// or "all_hidden" when all entries were filtered by ignore rules.
func RenderDirectoryPage(w io.Writer, displayPath string, parentHref *string, items []DirEntry, breadcrumbs []BreadcrumbSegment, emptyReason string) error {
	data := directoryData{
		DisplayPath: displayPath,
		HasParent:   parentHref != nil,
		ParentHref:  "",
		Breadcrumbs: breadcrumbs,
		Items:       items,
		EmptyReason: emptyReason,
		Version:     Version,
		WatchMode:   WatchMode,
	}
	if parentHref != nil {
		data.ParentHref = *parentHref
	}
	return templates.ExecuteTemplate(w, "directory.html", data)
}
