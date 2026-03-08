package render

import (
	"html/template"

	"github.com/driangle/viewmd/apps/cli/internal/frontmatter"
)

// BreadcrumbSegment represents one clickable segment in a file path breadcrumb.
type BreadcrumbSegment struct {
	Name string
	Href string
}

// markdownData holds template data for the markdown page.
type markdownData struct {
	FileName        string
	BaseURL         string
	ParentHref      string
	Breadcrumbs     []BreadcrumbSegment
	FrontmatterRows []frontmatter.KeyValue
	BodyHTML        template.HTML
	RawContent      string
	Version         string
}

// textData holds template data for the text file page.
type textData struct {
	FileName    string
	ParentHref  string
	Breadcrumbs []BreadcrumbSegment
	Content     template.HTML
	RawContent  string
	Language    string
	Version     string
}

// DirEntry represents a single item in a directory listing.
type DirEntry struct {
	Name  string
	Href  string
	IsDir bool
}

// CSSClass returns "dir" for directories and "file" for regular files.
func (e DirEntry) CSSClass() string {
	if e.IsDir {
		return "dir"
	}
	return "file"
}

// Label returns the display name, appending "/" for directories.
func (e DirEntry) Label() string {
	if e.IsDir {
		return e.Name + "/"
	}
	return e.Name
}

// unsupportedData holds template data for the unsupported file page.
type unsupportedData struct {
	FileName     string
	FileType     string
	FileSize     string
	DownloadHref string
	ParentHref   string
	Breadcrumbs  []BreadcrumbSegment
	Version      string
}

// directoryData holds template data for the directory listing page.
type directoryData struct {
	DisplayPath string
	HasParent   bool
	ParentHref  string
	Breadcrumbs []BreadcrumbSegment
	Items       []DirEntry
	Version     string
}
