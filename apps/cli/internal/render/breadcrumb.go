package render

import "strings"

// BuildBreadcrumbs creates breadcrumb segments from a URL request path.
// rootLabel is the display name for the root segment (typically the absolute path).
// For example, BuildBreadcrumbs("docs/api/file.md", "/home/user/project") produces:
//
//	[{Name: "/home/user/project", Href: "/"}, {Name: "docs", Href: "/docs/"}, {Name: "api", Href: "/docs/api/"}, {Name: "file.md", Href: ""}]
//
// The last segment (the current file) has an empty Href since it's not a link.
func BuildBreadcrumbs(reqPath string, rootLabel string) []BreadcrumbSegment {
	segments := []BreadcrumbSegment{
		{Name: rootLabel, Href: "/"},
	}

	reqPath = strings.TrimPrefix(reqPath, "/")
	reqPath = strings.TrimSuffix(reqPath, "/")
	if reqPath == "" || reqPath == "." {
		// Root directory — "root" is the current page
		segments[0].Href = ""
		return segments
	}

	parts := strings.Split(reqPath, "/")
	for i, part := range parts {
		isLast := i == len(parts)-1
		href := "/" + strings.Join(parts[:i+1], "/")
		if !isLast {
			href += "/"
		}
		if isLast {
			href = "" // current page, no link
		}
		segments = append(segments, BreadcrumbSegment{
			Name: part,
			Href: href,
		})
	}

	return segments
}
