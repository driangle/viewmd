package export

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteZip(t *testing.T) {
	// Set up temp directory with mixed content
	dir := t.TempDir()
	writeFile(t, dir, "readme.md", "# Hello\nWorld")
	writeFile(t, dir, "notes.markdown", "## Notes")
	writeFile(t, dir, "image.png", "fake-png-data")
	writeFile(t, dir, "sub/nested.md", "# Nested")
	writeFile(t, dir, "sub/data.txt", "plain text")

	var buf bytes.Buffer
	if err := WriteZip(&buf, dir, nil, false, "light"); err != nil {
		t.Fatalf("WriteZip: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("zip.NewReader: %v", err)
	}

	files := make(map[string]string)
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("open %s: %v", f.Name, err)
		}
		var content bytes.Buffer
		content.ReadFrom(rc)
		rc.Close()
		files[f.Name] = content.String()
	}

	tests := []struct {
		name      string
		wantExist bool
		checkFunc func(content string) bool
	}{
		{"readme.html", true, func(c string) bool {
			return strings.Contains(c, "<h1>Hello</h1>") &&
				strings.Contains(c, "<!DOCTYPE html>") &&
				strings.Contains(c, "--color-text:")
		}},
		{"notes.html", true, func(c string) bool { return strings.Contains(c, "<h2>Notes</h2>") }},
		{"readme.md", false, nil},
		{"notes.markdown", false, nil},
		{"image.png", true, func(c string) bool { return c == "fake-png-data" }},
		{"sub/nested.html", true, func(c string) bool { return strings.Contains(c, "<h1>Nested</h1>") }},
		{"sub/nested.md", false, nil},
		{"sub/data.txt", true, func(c string) bool { return c == "plain text" }},
	}

	for _, tt := range tests {
		content, exists := files[tt.name]
		if exists != tt.wantExist {
			t.Errorf("%s: exists=%v, want %v (files: %v)", tt.name, exists, tt.wantExist, fileNames(zr))
			continue
		}
		if tt.checkFunc != nil && !tt.checkFunc(content) {
			t.Errorf("%s: unexpected content: %q", tt.name, content)
		}
	}
}

func TestWriteZipFrontmatter(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "doc.md", "---\ntitle: My Doc\nauthor: Jane\n---\n# Content")

	var buf bytes.Buffer
	if err := WriteZip(&buf, dir, nil, false, "light"); err != nil {
		t.Fatalf("WriteZip: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("zip.NewReader: %v", err)
	}

	for _, f := range zr.File {
		if f.Name != "doc.html" {
			continue
		}
		rc, _ := f.Open()
		var content bytes.Buffer
		content.ReadFrom(rc)
		rc.Close()
		c := content.String()
		if !strings.Contains(c, "My Doc") {
			t.Error("frontmatter title not in exported HTML")
		}
		if !strings.Contains(c, "fm-key") {
			t.Error("frontmatter styling class not in exported HTML")
		}
		if !strings.Contains(c, "<h1>Content</h1>") {
			t.Error("markdown body not in exported HTML")
		}
		return
	}
	t.Error("doc.html not found in ZIP")
}

func TestWriteZipWithIgnore(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "keep.md", "# Keep")
	writeFile(t, dir, "secret.env", "PASSWORD=123")

	isIgnored := func(name, _ string) bool {
		return filepath.Ext(name) == ".env"
	}

	var buf bytes.Buffer
	if err := WriteZip(&buf, dir, isIgnored, false, "light"); err != nil {
		t.Fatalf("WriteZip: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("zip.NewReader: %v", err)
	}

	names := fileNames(zr)
	if len(names) != 1 || names[0] != "keep.html" {
		t.Errorf("expected [keep.html], got %v", names)
	}
}

func TestWriteZipMarkdownOnly(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "readme.md", "# Hello")
	writeFile(t, dir, "data.txt", "plain text")
	writeFile(t, dir, "image.png", "fake-png")
	writeFile(t, dir, "sub/nested.md", "# Nested")
	writeFile(t, dir, "sub/config.json", `{}`)

	var buf bytes.Buffer
	if err := WriteZip(&buf, dir, nil, true, "light"); err != nil {
		t.Fatalf("WriteZip: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("zip.NewReader: %v", err)
	}

	names := fileNames(zr)
	want := map[string]bool{"readme.html": true, "sub/nested.html": true}
	unwant := []string{"data.txt", "image.png", "sub/config.json"}

	for w := range want {
		found := false
		for _, n := range names {
			if n == w {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing %q in ZIP (got %v)", w, names)
		}
	}
	for _, u := range unwant {
		for _, n := range names {
			if n == u {
				t.Errorf("unexpected %q in ZIP with markdownOnly=true", u)
			}
		}
	}
	if len(names) != 2 {
		t.Errorf("expected 2 files, got %d: %v", len(names), names)
	}
}

func TestReplaceExt(t *testing.T) {
	tests := []struct {
		path, newExt, want string
	}{
		{"readme.md", ".html", "readme.html"},
		{"sub/doc.markdown", ".html", "sub/doc.html"},
		{"noext", ".html", "noext.html"},
	}
	for _, tt := range tests {
		got := replaceExt(tt.path, tt.newExt)
		if got != tt.want {
			t.Errorf("replaceExt(%q, %q) = %q, want %q", tt.path, tt.newExt, got, tt.want)
		}
	}
}

func writeFile(t *testing.T, base, relPath, content string) {
	t.Helper()
	full := filepath.Join(base, relPath)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fileNames(zr *zip.Reader) []string {
	var names []string
	for _, f := range zr.File {
		names = append(names, f.Name)
	}
	return names
}
