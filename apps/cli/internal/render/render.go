// Package render provides HTML rendering for markdown content,
// text files, and directory listings.
package render

import (
	"embed"
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
func RenderMarkdownPage(w io.Writer, fileName string, frontmatter map[string]string, bodyHTML string, baseURL string, parentHref string, rawContent string) error {
	var rows []FrontmatterRow
	for k, v := range frontmatter {
		rows = append(rows, FrontmatterRow{Key: k, Value: v})
	}

	data := markdownData{
		FileName:        fileName,
		BaseURL:         baseURL,
		ParentHref:      parentHref,
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
func RenderTextPage(w io.Writer, fileName string, escapedContent string, parentHref string, rawContent string) error {
	data := textData{
		FileName:   fileName,
		ParentHref: parentHref,
		Content:    template.HTML(escapedContent),
		RawContent: rawContent,
		Version:    Version,
	}
	return templates.ExecuteTemplate(w, "text.html", data)
}

// RenderDirectoryPage writes a full HTML page for a directory listing.
// parentHref is nil for the root directory.
func RenderDirectoryPage(w io.Writer, displayPath string, parentHref *string, items []DirEntry) error {
	data := directoryData{
		DisplayPath: displayPath,
		HasParent:   parentHref != nil,
		ParentHref:  "",
		Items:       items,
		Version:     Version,
	}
	if parentHref != nil {
		data.ParentHref = *parentHref
	}
	return templates.ExecuteTemplate(w, "directory.html", data)
}
