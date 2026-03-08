package render

import (
	"testing"
)

func TestBuildBreadcrumbs(t *testing.T) {
	tests := []struct {
		name    string
		reqPath string
		want    []BreadcrumbSegment
	}{
		{
			name:    "empty path (root dir)",
			reqPath: "",
			want:    []BreadcrumbSegment{{Name: "root", Href: ""}},
		},
		{
			name:    "root file",
			reqPath: "file.md",
			want: []BreadcrumbSegment{
				{Name: "root", Href: "/"},
				{Name: "file.md", Href: ""},
			},
		},
		{
			name:    "nested file",
			reqPath: "docs/file.md",
			want: []BreadcrumbSegment{
				{Name: "root", Href: "/"},
				{Name: "docs", Href: "/docs/"},
				{Name: "file.md", Href: ""},
			},
		},
		{
			name:    "deeply nested file",
			reqPath: "a/b/c/file.txt",
			want: []BreadcrumbSegment{
				{Name: "root", Href: "/"},
				{Name: "a", Href: "/a/"},
				{Name: "b", Href: "/a/b/"},
				{Name: "c", Href: "/a/b/c/"},
				{Name: "file.txt", Href: ""},
			},
		},
		{
			name:    "leading slash stripped",
			reqPath: "/docs/file.md",
			want: []BreadcrumbSegment{
				{Name: "root", Href: "/"},
				{Name: "docs", Href: "/docs/"},
				{Name: "file.md", Href: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildBreadcrumbs(tt.reqPath, "root")
			if len(got) != len(tt.want) {
				t.Fatalf("BuildBreadcrumbs(%q) returned %d segments, want %d", tt.reqPath, len(got), len(tt.want))
			}
			for i, seg := range got {
				if seg.Name != tt.want[i].Name || seg.Href != tt.want[i].Href {
					t.Errorf("segment[%d] = {%q, %q}, want {%q, %q}", i, seg.Name, seg.Href, tt.want[i].Name, tt.want[i].Href)
				}
			}
		})
	}
}
