package render

import "html/template"

// FrontmatterRow holds a single key-value pair for the frontmatter table.
type FrontmatterRow struct {
	Key   string
	Value string
}

// markdownData holds template data for the markdown page.
type markdownData struct {
	FileName        string
	BaseURL         string
	ParentHref      string
	Frontmatter     map[string]string
	FrontmatterRows []FrontmatterRow
	BodyHTML        template.HTML
	RawContent      string
	Version         string
}

// textData holds template data for the text file page.
type textData struct {
	FileName   string
	ParentHref string
	Content    template.HTML
	RawContent string
	Version    string
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

// directoryData holds template data for the directory listing page.
type directoryData struct {
	DisplayPath string
	HasParent   bool
	ParentHref  string
	Items       []DirEntry
	Version     string
}
