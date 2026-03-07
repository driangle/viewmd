package render_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driangle/viewmd/apps/cli/internal/render"
)

func TestMarkdownPageWithFrontmatter(t *testing.T) {
	var buf bytes.Buffer
	fm := map[string]string{"title": "Hello", "author": "Alice"}
	err := render.RenderMarkdownPage(&buf, "test.md", fm, "<p>Body</p>", "/", "/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, `class="frontmatter"`) {
		t.Error("expected .frontmatter class in output")
	}
	if !strings.Contains(out, `class="fm-key"`) {
		t.Error("expected .fm-key class in output")
	}
	if !strings.Contains(out, "<p>Body</p>") {
		t.Error("expected body HTML in output")
	}
}

func TestMarkdownPageWithoutFrontmatter(t *testing.T) {
	var buf bytes.Buffer
	err := render.RenderMarkdownPage(&buf, "test.md", nil, "<p>Body</p>", "/", "/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if strings.Contains(out, `<div class="frontmatter"`) {
		t.Error("expected no frontmatter div when frontmatter is nil")
	}
	if !strings.Contains(out, "<p>Body</p>") {
		t.Error("expected body HTML in output")
	}
}

func TestMarkdownPageBaseHref(t *testing.T) {
	var buf bytes.Buffer
	err := render.RenderMarkdownPage(&buf, "doc.md", nil, "<p>hi</p>", "/docs/", "/docs/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, `<base href="/docs/">`) {
		t.Error("expected base href to be set to /docs/")
	}
}

func TestMarkdownPageEscapesFileName(t *testing.T) {
	var buf bytes.Buffer
	err := render.RenderMarkdownPage(&buf, "<script>alert(1)</script>", nil, "<p>ok</p>", "/", "/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if strings.Contains(out, "<script>alert(1)</script>") {
		t.Error("expected file name to be HTML-escaped")
	}
	if !strings.Contains(out, "&lt;script&gt;") {
		t.Error("expected escaped script tag in title")
	}
}

func TestMarkdownPageEscapesFrontmatterValues(t *testing.T) {
	var buf bytes.Buffer
	fm := map[string]string{"key": "<b>bold</b>"}
	err := render.RenderMarkdownPage(&buf, "test.md", fm, "<p>ok</p>", "/", "/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if strings.Contains(out, "<b>bold</b>") {
		t.Error("expected frontmatter values to be HTML-escaped")
	}
	if !strings.Contains(out, "&lt;b&gt;bold&lt;/b&gt;") {
		t.Error("expected escaped bold tag in frontmatter")
	}
}

func TestTextPage(t *testing.T) {
	var buf bytes.Buffer
	err := render.RenderTextPage(&buf, "main.go", "package main\n", "/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, `<div class="header">main.go</div>`) {
		t.Error("expected file name in header div")
	}
	if !strings.Contains(out, "<pre>package main\n</pre>") {
		t.Error("expected content in pre block")
	}
}

func TestTextPageEscapesFileName(t *testing.T) {
	var buf bytes.Buffer
	err := render.RenderTextPage(&buf, "<img src=x>", "content", "/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if strings.Contains(out, "<img src=x>") {
		t.Error("expected file name to be HTML-escaped in text page")
	}
}

func TestDirectoryPageWithParent(t *testing.T) {
	var buf bytes.Buffer
	parent := "parent"
	items := []render.DirEntry{
		{Name: "subdir", Href: "path/subdir", IsDir: true},
		{Name: "file.txt", Href: "path/file.txt", IsDir: false},
	}
	err := render.RenderDirectoryPage(&buf, "path", &parent, items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, `<a href="/parent" class="dir">..</a>`) {
		t.Error("expected parent link with dir class")
	}
	if !strings.Contains(out, `<a href="/path/subdir" class="dir">subdir/</a>`) {
		t.Error("expected directory entry with trailing slash and dir class")
	}
	if !strings.Contains(out, `<a href="/path/file.txt" class="file">file.txt</a>`) {
		t.Error("expected file entry with file class")
	}
	if !strings.Contains(out, "<h1>Directory: /path</h1>") {
		t.Error("expected display path in h1")
	}
}

func TestMarkdownPageParentNav(t *testing.T) {
	tests := []struct {
		name       string
		parentHref string
		wantLink   string
	}{
		{"root file", "/", `<a href="/">..</a>`},
		{"nested file", "/docs/", `<a href="/docs/">..</a>`},
		{"deeply nested", "/a/b/", `<a href="/a/b/">..</a>`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := render.RenderMarkdownPage(&buf, "f.md", nil, "<p>ok</p>", "/", tt.parentHref)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			out := buf.String()

			if !strings.Contains(out, `class="parent-nav"`) {
				t.Error("expected parent-nav class")
			}
			if !strings.Contains(out, tt.wantLink) {
				t.Errorf("expected parent link %q in output", tt.wantLink)
			}
		})
	}
}

func TestTextPageParentNav(t *testing.T) {
	tests := []struct {
		name       string
		parentHref string
		wantLink   string
	}{
		{"root file", "/", `<a href="/">..</a>`},
		{"nested file", "/src/", `<a href="/src/">..</a>`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := render.RenderTextPage(&buf, "main.go", "code", tt.parentHref)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			out := buf.String()

			if !strings.Contains(out, `class="parent-nav"`) {
				t.Error("expected parent-nav class")
			}
			if !strings.Contains(out, tt.wantLink) {
				t.Errorf("expected parent link %q in output", tt.wantLink)
			}
		})
	}
}

func TestDirectoryPageWithoutParent(t *testing.T) {
	var buf bytes.Buffer
	items := []render.DirEntry{
		{Name: "readme.md", Href: "readme.md", IsDir: false},
	}
	err := render.RenderDirectoryPage(&buf, "", nil, items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if strings.Contains(out, "..") {
		t.Error("expected no parent link when parentHref is nil")
	}
}
