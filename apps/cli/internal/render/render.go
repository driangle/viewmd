// Package render provides HTML rendering for markdown content,
// text files, and directory listings.
package render

import (
	"embed"
	"fmt"
	"html/template"
	"io"
)

// Version is set by the main package at startup.
var Version string

//go:embed templates/*.html
var templateFS embed.FS

var templates = template.Must(
	template.ParseFS(templateFS, "templates/*.html"),
)

// RenderMarkdownPage writes a full HTML page for rendered markdown content.
// frontmatter may be nil if no frontmatter was parsed.
// bodyHTML is already-rendered HTML from the markdown converter.
func RenderMarkdownPage(w io.Writer, fileName string, frontmatter map[string]string, bodyHTML string, baseURL string, parentHref string, rawContent string, breadcrumbs []BreadcrumbSegment) error {
	var rows []FrontmatterRow
	for k, v := range frontmatter {
		rows = append(rows, FrontmatterRow{Key: k, Value: v})
	}

	data := markdownData{
		FileName:        fileName,
		BaseURL:         baseURL,
		ParentHref:      parentHref,
		Breadcrumbs:     breadcrumbs,
		Frontmatter:     frontmatter,
		FrontmatterRows: rows,
		BodyHTML:        template.HTML(bodyHTML),
		RawContent:      rawContent,
		Version:         Version,
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

// RenderDirectoryPage writes a full HTML page for a directory listing.
// parentHref is nil for the root directory.
func RenderDirectoryPage(w io.Writer, displayPath string, parentHref *string, items []DirEntry, breadcrumbs []BreadcrumbSegment) error {
	data := directoryData{
		DisplayPath: displayPath,
		HasParent:   parentHref != nil,
		ParentHref:  "",
		Breadcrumbs: breadcrumbs,
		Items:       items,
		Version:     Version,
	}
	if parentHref != nil {
		data.ParentHref = *parentHref
	}
	return templates.ExecuteTemplate(w, "directory.html", data)
}
