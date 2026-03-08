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
	err := render.RenderMarkdownPage(&buf, "test.md", fm, "<p>Body</p>", "/", "/", "# Body", nil)
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
	err := render.RenderMarkdownPage(&buf, "test.md", nil, "<p>Body</p>", "/", "/", "Body", nil)
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
	err := render.RenderMarkdownPage(&buf, "doc.md", nil, "<p>hi</p>", "/docs/", "/docs/", "hi", nil)
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
	err := render.RenderMarkdownPage(&buf, "<script>alert(1)</script>", nil, "<p>ok</p>", "/", "/", "ok", nil)
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
	err := render.RenderMarkdownPage(&buf, "test.md", fm, "<p>ok</p>", "/", "/", "ok", nil)
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
	err := render.RenderTextPage(&buf, "main.go", "package main\n", "/", "package main\n", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, `<span>main.go</span>`) {
		t.Error("expected file name in header span")
	}
	if !strings.Contains(out, "<pre>package main\n</pre>") {
		t.Error("expected content in pre block")
	}
}

func TestTextPageEscapesFileName(t *testing.T) {
	var buf bytes.Buffer
	err := render.RenderTextPage(&buf, "<img src=x>", "content", "/", "content", nil)
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
	breadcrumbs := render.BuildBreadcrumbs("path")
	err := render.RenderDirectoryPage(&buf, "path", &parent, items, breadcrumbs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, `class="breadcrumb"`) {
		t.Error("expected breadcrumb nav")
	}
	if !strings.Contains(out, `<a href="/">root</a>`) {
		t.Error("expected root breadcrumb link")
	}
	if !strings.Contains(out, `<a href="/path/subdir" class="dir">subdir/</a>`) {
		t.Error("expected directory entry with trailing slash and dir class")
	}
	if !strings.Contains(out, `<a href="/path/file.txt" class="file">file.txt</a>`) {
		t.Error("expected file entry with file class")
	}
}

func TestMarkdownPageBreadcrumb(t *testing.T) {
	var buf bytes.Buffer
	breadcrumbs := render.BuildBreadcrumbs("docs/f.md")
	err := render.RenderMarkdownPage(&buf, "f.md", nil, "<p>ok</p>", "/", "/docs/", "ok", breadcrumbs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, `class="breadcrumb"`) {
		t.Error("expected breadcrumb class")
	}
	if !strings.Contains(out, `<a href="/">root</a>`) {
		t.Error("expected root breadcrumb link")
	}
	if !strings.Contains(out, `<a href="/docs/">docs</a>`) {
		t.Error("expected docs breadcrumb link")
	}
}

func TestTextPageBreadcrumb(t *testing.T) {
	var buf bytes.Buffer
	breadcrumbs := render.BuildBreadcrumbs("src/main.go")
	err := render.RenderTextPage(&buf, "main.go", "code", "/src/", "code", breadcrumbs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, `class="breadcrumb"`) {
		t.Error("expected breadcrumb class")
	}
	if !strings.Contains(out, `<a href="/">root</a>`) {
		t.Error("expected root breadcrumb link")
	}
}

func TestDirectoryPageWithoutParent(t *testing.T) {
	var buf bytes.Buffer
	items := []render.DirEntry{
		{Name: "readme.md", Href: "readme.md", IsDir: false},
	}
	breadcrumbs := render.BuildBreadcrumbs("")
	err := render.RenderDirectoryPage(&buf, "", nil, items, breadcrumbs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, `<span class="current">root</span>`) {
		t.Error("expected root as current in breadcrumb for root directory")
	}
}
