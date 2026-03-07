// Package render provides HTML rendering for markdown content,
// text files, and directory listings.
package render

import (
	"embed"
	"html/template"
	"io"
)

//go:embed templates/*.html
var templateFS embed.FS

var templates = template.Must(
	template.ParseFS(templateFS, "templates/*.html"),
)

// RenderMarkdownPage writes a full HTML page for rendered markdown content.
// frontmatter may be nil if no frontmatter was parsed.
// bodyHTML is already-rendered HTML from the markdown converter.
func RenderMarkdownPage(w io.Writer, fileName string, frontmatter map[string]string, bodyHTML string, baseURL string) error {
	var rows []FrontmatterRow
	for k, v := range frontmatter {
		rows = append(rows, FrontmatterRow{Key: k, Value: v})
	}

	data := markdownData{
		FileName:        fileName,
		BaseURL:         baseURL,
		Frontmatter:     frontmatter,
		FrontmatterRows: rows,
		BodyHTML:        template.HTML(bodyHTML),
	}
	return templates.ExecuteTemplate(w, "markdown.html", data)
}

// RenderTextPage writes a full HTML page for a plain text file.
// escapedContent should already be HTML-escaped.
func RenderTextPage(w io.Writer, fileName string, escapedContent string) error {
	data := textData{
		FileName: fileName,
		Content:  template.HTML(escapedContent),
	}
	return templates.ExecuteTemplate(w, "text.html", data)
}

// RenderDirectoryPage writes a full HTML page for a directory listing.
// parentHref is nil for the root directory.
func RenderDirectoryPage(w io.Writer, displayPath string, parentHref *string, items []DirEntry) error {
	var parent string
	if parentHref != nil {
		parent = *parentHref
	}

	data := directoryData{
		DisplayPath: displayPath,
		ParentHref:  parent,
		Items:       items,
	}
	return templates.ExecuteTemplate(w, "directory.html", data)
}
