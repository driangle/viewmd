package markdown

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "heading",
			input:    "# Hello",
			contains: "<h1>Hello</h1>",
		},
		{
			name:     "fenced code block",
			input:    "```\nfoo()\n```",
			contains: "<pre><code>",
		},
		{
			name:     "table",
			input:    "| A | B |\n|---|---|\n| 1 | 2 |",
			contains: "<table>",
		},
		{
			name:     "newline to br",
			input:    "line1\nline2",
			contains: "<br",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.input)
			if err != nil {
				t.Fatalf("Convert() error: %v", err)
			}
			if !strings.Contains(got, tt.contains) {
				t.Errorf("Convert(%q) = %q, want substring %q", tt.input, got, tt.contains)
			}
		})
	}
}

func TestServeMarkdown(t *testing.T) {
	dir := t.TempDir()
	mdFile := filepath.Join(dir, "test.md")
	content := "---\ntitle: Hello\n---\n# Welcome\n\nSome **bold** text."
	if err := os.WriteFile(mdFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	ServeMarkdown(rec, mdFile, "/", "/", nil)

	if rec.Code != 200 {
		t.Errorf("status = %d, want 200", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}

	body := rec.Body.String()
	for _, want := range []string{"<h1>Welcome</h1>", "<strong>bold</strong>", "Hello"} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q", want)
		}
	}
}

func TestServeMarkdown_FileNotFound(t *testing.T) {
	rec := httptest.NewRecorder()
	ServeMarkdown(rec, "/nonexistent/path.md", "/", "/", nil)

	if rec.Code != 500 {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

