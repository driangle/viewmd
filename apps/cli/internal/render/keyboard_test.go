package render_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driangle/viewmd/apps/cli/internal/render"
)

func TestDirectoryPageKeyboardNavJS(t *testing.T) {
	var buf bytes.Buffer
	parent := "parent"
	items := []render.DirEntry{
		{Name: "file.txt", Href: "file.txt", IsDir: false},
	}
	err := render.RenderDirectoryPage(&buf, "path", &parent, items, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	checks := []struct {
		name    string
		content string
	}{
		{"file list id", `id="file-list"`},
		{"keydown listener", `addEventListener('keydown'`},
		{"ArrowDown handler", `e.key === 'ArrowDown'`},
		{"ArrowUp handler", `e.key === 'ArrowUp'`},
		{"Enter handler", `e.key === 'Enter'`},
		{"ArrowRight handler", `e.key === 'ArrowRight'`},
		{"ArrowLeft handler", `e.key === 'ArrowLeft'`},
		{"Backspace handler", `e.key === 'Backspace'`},
		{"kb-active class", `kb-active`},
	}
	for _, c := range checks {
		if !strings.Contains(out, c.content) {
			t.Errorf("expected %s (%q) in directory page output", c.name, c.content)
		}
	}
}

func TestDirectoryPageKeyboardHint(t *testing.T) {
	var buf bytes.Buffer
	items := []render.DirEntry{
		{Name: "file.txt", Href: "file.txt", IsDir: false},
	}
	err := render.RenderDirectoryPage(&buf, "", nil, items, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, `class="kb-hint"`) {
		t.Error("expected keyboard hint element with kb-hint class")
	}
	if !strings.Contains(out, "navigate") {
		t.Error("expected 'navigate' in keyboard hint text")
	}
	if !strings.Contains(out, "open") {
		t.Error("expected 'open' in keyboard hint text")
	}
	if !strings.Contains(out, "back") {
		t.Error("expected 'back' in keyboard hint text")
	}
}

func TestDirectoryPageHighlightCSS(t *testing.T) {
	var buf bytes.Buffer
	items := []render.DirEntry{
		{Name: "a.txt", Href: "a.txt", IsDir: false},
	}
	err := render.RenderDirectoryPage(&buf, "", nil, items, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "li.kb-active") {
		t.Error("expected CSS rule for li.kb-active highlight style")
	}
}

func TestDirectoryPageBreadcrumbParentHref(t *testing.T) {
	var buf bytes.Buffer
	parent := ""
	items := []render.DirEntry{}
	breadcrumbs := render.BuildBreadcrumbs("sub")
	err := render.RenderDirectoryPage(&buf, "sub", &parent, items, breadcrumbs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "data-parent-href") {
		t.Error("expected data-parent-href attribute on breadcrumb")
	}
}

func TestMarkdownPageBackNavJS(t *testing.T) {
	var buf bytes.Buffer
	err := render.RenderMarkdownPage(&buf, "test.md", nil, "<p>hi</p>", "/", "/", "hi", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, `addEventListener('keydown'`) {
		t.Error("expected keydown listener on markdown page")
	}
	if !strings.Contains(out, `e.key === 'ArrowLeft'`) {
		t.Error("expected ArrowLeft handler on markdown page")
	}
	if !strings.Contains(out, `e.key === 'Backspace'`) {
		t.Error("expected Backspace handler on markdown page")
	}
	if !strings.Contains(out, `.breadcrumb`) {
		t.Error("expected breadcrumb selector in back-nav script")
	}
}

func TestTextPageBackNavJS(t *testing.T) {
	var buf bytes.Buffer
	err := render.RenderTextPage(&buf, "main.go", "code", "/", "code", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, `addEventListener('keydown'`) {
		t.Error("expected keydown listener on text page")
	}
	if !strings.Contains(out, `e.key === 'ArrowLeft'`) {
		t.Error("expected ArrowLeft handler on text page")
	}
	if !strings.Contains(out, `e.key === 'Backspace'`) {
		t.Error("expected Backspace handler on text page")
	}
	if !strings.Contains(out, `.breadcrumb`) {
		t.Error("expected breadcrumb selector in back-nav script")
	}
}
