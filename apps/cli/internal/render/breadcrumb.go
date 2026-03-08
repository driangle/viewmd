package render

import "strings"

// BuildBreadcrumbs creates breadcrumb segments from a URL request path.
// For example, "docs/api/file.md" produces:
//
//	[{Name: "root", Href: "/"}, {Name: "docs", Href: "/docs/"}, {Name: "api", Href: "/docs/api/"}, {Name: "file.md", Href: ""}]
//
// The last segment (the current file) has an empty Href since it's not a link.
func BuildBreadcrumbs(reqPath string) []BreadcrumbSegment {
	segments := []BreadcrumbSegment{
		{Name: "root", Href: "/"},
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
