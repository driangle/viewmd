package handler_test

import (
	"archive/zip"
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/driangle/viewmd/apps/cli/internal/handler"
)

func TestZipExportDirectory(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "readme.md", "# Hello")
	writeFile(t, dir, "data.txt", "plain text")
	os.MkdirAll(filepath.Join(dir, "sub"), 0o755)
	writeFile(t, dir, "sub/nested.md", "# Nested")

	h := handler.New(dir)
	h.ShowAll = true

	rec := request(h, "/?export=zip")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/zip" {
		t.Errorf("Content-Type = %q, want application/zip", ct)
	}
	if cd := rec.Header().Get("Content-Disposition"); !strings.Contains(cd, ".zip") {
		t.Errorf("Content-Disposition = %q, want to contain .zip", cd)
	}

	zr, err := zip.NewReader(bytes.NewReader(rec.Body.Bytes()), int64(rec.Body.Len()))
	if err != nil {
		t.Fatalf("zip.NewReader: %v", err)
	}

	names := make(map[string]bool)
	for _, f := range zr.File {
		names[f.Name] = true
	}

	want := []string{"readme.html", "data.txt", "sub/nested.html"}
	for _, w := range want {
		if !names[w] {
			t.Errorf("missing %q in ZIP (got %v)", w, names)
		}
	}

	// Markdown originals should not be present
	unwant := []string{"readme.md", "sub/nested.md"}
	for _, u := range unwant {
		if names[u] {
			t.Errorf("unexpected %q in ZIP", u)
		}
	}
}

func TestZipExportSubdirectory(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "docs"), 0o755)
	writeFile(t, dir, "docs/guide.md", "# Guide")
	writeFile(t, dir, "docs/config.json", `{"key":"val"}`)

	h := handler.New(dir)
	h.ShowAll = true

	rec := request(h, "/docs?export=zip")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if cd := rec.Header().Get("Content-Disposition"); !strings.Contains(cd, "docs.zip") {
		t.Errorf("Content-Disposition = %q, want docs.zip", cd)
	}

	zr, err := zip.NewReader(bytes.NewReader(rec.Body.Bytes()), int64(rec.Body.Len()))
	if err != nil {
		t.Fatalf("zip.NewReader: %v", err)
	}

	names := make(map[string]bool)
	for _, f := range zr.File {
		names[f.Name] = true
	}

	if !names["guide.html"] {
		t.Error("missing guide.html")
	}
	if !names["config.json"] {
		t.Error("missing config.json")
	}
	if names["guide.md"] {
		t.Error("unexpected guide.md (should be converted)")
	}
}

func TestZipExportRespectsIgnorePatterns(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "keep.md", "# Keep")
	writeFile(t, dir, "secret.env", "PASSWORD=123")

	h := handler.New(dir)
	h.ShowAll = true
	h.IgnorePatterns = []string{"*.env"}

	rec := request(h, "/?export=zip")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	zr, err := zip.NewReader(bytes.NewReader(rec.Body.Bytes()), int64(rec.Body.Len()))
	if err != nil {
		t.Fatalf("zip.NewReader: %v", err)
	}

	for _, f := range zr.File {
		if strings.HasSuffix(f.Name, ".env") {
			t.Errorf("ignored file %q should not be in ZIP", f.Name)
		}
	}
}

func TestHTMLExportSingleFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "doc.md", "---\ntitle: Test\n---\n# Hello World")

	h := handler.New(dir)
	h.ShowAll = true

	rec := request(h, "/doc.md?export=html")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
	if cd := rec.Header().Get("Content-Disposition"); !strings.Contains(cd, "doc.html") {
		t.Errorf("Content-Disposition = %q, want to contain doc.html", cd)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("missing DOCTYPE in exported HTML")
	}
	if !strings.Contains(body, "<h1>Hello World</h1>") {
		t.Error("missing rendered markdown in exported HTML")
	}
	if !strings.Contains(body, "--color-text:") {
		t.Error("missing CSS styling in exported HTML")
	}
	if !strings.Contains(body, "Test") {
		t.Error("missing frontmatter in exported HTML")
	}
}

func TestHTMLExportNonMarkdownIgnored(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "script.py", "print('hello')")

	h := handler.New(dir)
	h.ShowAll = true

	rec := request(h, "/script.py?export=html")

	// Should render normally, not as an export download
	if cd := rec.Header().Get("Content-Disposition"); cd != "" {
		t.Errorf("non-markdown file should not get export disposition, got %q", cd)
	}
}
