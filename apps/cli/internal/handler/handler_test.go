package handler_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/driangle/viewmd/apps/cli/internal/handler"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// Create test files
	writeFile(t, dir, "README.md", "# Hello\nWorld")
	writeFile(t, dir, "doc.md", "---\ntitle: Test\n---\n# Doc")
	writeFile(t, dir, "notes.markdown", "# Notes")
	writeFile(t, dir, "script.py", "print('hello')")
	writeFile(t, dir, ".gitignore", "*.o")

	// Create a subdirectory with files
	os.MkdirAll(filepath.Join(dir, "sub"), 0o755)
	writeFile(t, dir, "sub/page.md", "# Sub Page")

	// Create a subdirectory with README
	os.MkdirAll(filepath.Join(dir, "docs"), 0o755)
	writeFile(t, dir, "docs/README.md", "# Docs Index")

	// Create a binary file
	os.WriteFile(filepath.Join(dir, "image.png"), []byte{0x89, 0x50, 0x4E, 0x47}, 0o644)

	return dir
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func request(h http.Handler, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", path, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestRootServesDirectoryListing(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Directory:") {
		t.Error("expected directory listing page")
	}
}

func TestMarkdownFileServed(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/doc.md")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Doc") {
		t.Error("expected rendered markdown content")
	}
	if !strings.Contains(body, "title") {
		t.Error("expected frontmatter in output")
	}
}

func TestMarkdownExtension(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/notes.markdown")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Notes") {
		t.Error("expected rendered markdown for .markdown extension")
	}
}

func TestTextFileServed(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/script.py")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "print(") {
		t.Error("expected text file content")
	}
}

func TestDotfileServedAsText(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/.gitignore")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "*.o") {
		t.Error("expected dotfile content")
	}
}

func TestDirectoryWithReadmeServesReadme(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/docs")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Docs Index") {
		t.Error("expected README.md content for directory with README")
	}
}

func TestDirectoryWithoutReadmeShowsListing(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/sub")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Directory:") {
		t.Error("expected directory listing")
	}
	if !strings.Contains(body, "page.md") {
		t.Error("expected page.md in listing")
	}
}

func TestMissingFileReturns404(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/nonexistent.md")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("got status %d, want 404", rec.Code)
	}
}

func TestDirectoryListingSortOrder(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "zebra.txt", "z")
	writeFile(t, dir, "alpha.txt", "a")
	os.MkdirAll(filepath.Join(dir, "beta_dir"), 0o755)
	os.MkdirAll(filepath.Join(dir, "alpha_dir"), 0o755)

	h := handler.New(dir)
	rec := request(h, "/")

	body := rec.Body.String()
	// Dirs should come before files
	alphaDirPos := strings.Index(body, "alpha_dir")
	betaDirPos := strings.Index(body, "beta_dir")
	alphaPos := strings.Index(body, "alpha.txt")
	zebraPos := strings.Index(body, "zebra.txt")

	if alphaDirPos > betaDirPos {
		t.Error("alpha_dir should come before beta_dir")
	}
	if betaDirPos > alphaPos {
		t.Error("directories should come before files")
	}
	if alphaPos > zebraPos {
		t.Error("alpha.txt should come before zebra.txt")
	}
}

func TestSubdirectoryListingHasParentLink(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/sub")

	body := rec.Body.String()
	if !strings.Contains(body, "..") {
		t.Error("subdirectory listing should have parent link")
	}
}

func TestRootListingNoParentLink(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/")

	body := rec.Body.String()
	// The parent link in the template is rendered as ".." with class="dir"
	if strings.Contains(body, `class="dir">..</a>`) {
		t.Error("root listing should not have parent link")
	}
}

func TestBinaryFileServed(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/image.png")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
}

func TestURLEncodedPath(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "my file.md", "# Spaced")

	h := handler.New(dir)
	rec := request(h, "/my%20file.md")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Spaced") {
		t.Error("expected content from URL-encoded path")
	}
}
