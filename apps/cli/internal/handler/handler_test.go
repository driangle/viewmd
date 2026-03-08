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
	writeFile(t, dir, "unsafe.txt", "<script>alert('xss')</script>")
	writeFile(t, dir, ".gitignore", "*.o")

	// Create a subdirectory with files
	os.MkdirAll(filepath.Join(dir, "sub"), 0o755)
	writeFile(t, dir, "sub/page.md", "# Sub Page")

	// Create a subdirectory with README
	os.MkdirAll(filepath.Join(dir, "docs"), 0o755)
	writeFile(t, dir, "docs/README.md", "# Docs Index")

	// Create a binary file
	os.WriteFile(filepath.Join(dir, "image.png"), []byte{0x89, 0x50, 0x4E, 0x47}, 0o644)
	os.WriteFile(filepath.Join(dir, "broken.txt"), []byte{0xff, 0xfe, 0xfd, 0x00, 0x01}, 0o644)

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
	if !strings.Contains(body, "Directory listing") {
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
	if !strings.Contains(body, `<span>script.py</span>`) {
		t.Error("expected filename shown in header")
	}
	if !strings.Contains(body, "<pre>") {
		t.Error("expected content wrapped in pre block")
	}
}

func TestTextFileContentEscaped(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/unsafe.txt")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if strings.Contains(body, "<script>alert('xss')</script>") {
		t.Error("expected raw script tag to be escaped")
	}
	if !strings.Contains(body, "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;") {
		t.Error("expected escaped script tag in output")
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

func TestDirectoryWithReadmeShowsListingByDefault(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/docs")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Directory listing") {
		t.Error("expected directory listing by default, not README content")
	}
}

func TestDirectoryWithReadmeAutoServeEnabled(t *testing.T) {
	h := handler.New(setupTestDir(t))
	h.AutoReadme = true
	rec := request(h, "/docs")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Docs Index") {
		t.Error("expected README.md content when AutoReadme is enabled")
	}
}

func TestDirectoryWithoutReadmeShowsListing(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/sub")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Directory listing") {
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
	writeFile(t, dir, "zebra.md", "# Z")
	writeFile(t, dir, "alpha.md", "# A")
	os.MkdirAll(filepath.Join(dir, "beta_dir"), 0o755)
	writeFile(t, dir, "beta_dir/x.md", "# X")
	os.MkdirAll(filepath.Join(dir, "alpha_dir"), 0o755)
	writeFile(t, dir, "alpha_dir/y.md", "# Y")

	h := handler.New(dir)
	rec := request(h, "/")

	body := rec.Body.String()
	// Dirs should come before files
	alphaDirPos := strings.Index(body, "alpha_dir")
	betaDirPos := strings.Index(body, "beta_dir")
	alphaPos := strings.Index(body, "alpha.md")
	zebraPos := strings.Index(body, "zebra.md")

	if alphaDirPos > betaDirPos {
		t.Error("alpha_dir should come before beta_dir")
	}
	if betaDirPos > alphaPos {
		t.Error("directories should come before files")
	}
	if alphaPos > zebraPos {
		t.Error("alpha.md should come before zebra.md")
	}
}

func TestSubdirectoryListingHasBreadcrumb(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/sub")

	body := rec.Body.String()
	if !strings.Contains(body, `class="breadcrumb"`) {
		t.Error("subdirectory listing should have breadcrumb")
	}
	if !strings.Contains(body, `<a href="/">root</a>`) {
		t.Error("subdirectory breadcrumb should have root link")
	}
}

func TestRootListingBreadcrumbShowsOnlyRoot(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/")

	body := rec.Body.String()
	if !strings.Contains(body, `class="breadcrumb"`) {
		t.Error("root listing should have breadcrumb")
	}
	// Root should show "root" as current (no link)
	if !strings.Contains(body, `<span class="current">root</span>`) {
		t.Error("root listing breadcrumb should show root as current")
	}
}

func TestBinaryFileServesUnsupportedPage(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/image.png")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("got content-type %q, want text/html", rec.Header().Get("Content-Type"))
	}
	body := rec.Body.String()
	if !strings.Contains(body, "No preview available") {
		t.Fatal("expected unsupported page with 'No preview available'")
	}
	if !strings.Contains(body, "?raw=1") {
		t.Fatal("expected download link with ?raw=1")
	}
}

func TestRawQueryParamServesBinaryDirect(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/image.png?raw=1")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	if strings.Contains(rec.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("raw=1 should not serve HTML, got %q", rec.Header().Get("Content-Type"))
	}
}

func TestInvalidUTF8TextServesUnsupportedPage(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/broken.txt")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("got content-type %q, want text/html", rec.Header().Get("Content-Type"))
	}
	body := rec.Body.String()
	if !strings.Contains(body, "No preview available") {
		t.Fatal("expected unsupported page for invalid UTF-8 text file")
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

func TestMarkdownFileHasBreadcrumb(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/doc.md")

	body := rec.Body.String()
	if !strings.Contains(body, `class="breadcrumb"`) {
		t.Error("markdown file should have breadcrumb element")
	}
	if !strings.Contains(body, `<a href="/">root</a>`) {
		t.Error("breadcrumb should have root link")
	}
}

func TestNestedMarkdownFileHasBreadcrumb(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/sub/page.md")

	body := rec.Body.String()
	if !strings.Contains(body, `<a href="/sub/">sub</a>`) {
		t.Error("nested file breadcrumb should have parent directory link")
	}
}

func TestTextFileHasBreadcrumb(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/script.py")

	body := rec.Body.String()
	if !strings.Contains(body, `class="breadcrumb"`) {
		t.Error("text file should have breadcrumb element")
	}
	if !strings.Contains(body, `<a href="/">root</a>`) {
		t.Error("breadcrumb should have root link")
	}
}

func TestReadmeAutoServeHasBreadcrumb(t *testing.T) {
	h := handler.New(setupTestDir(t))
	h.AutoReadme = true
	rec := request(h, "/docs")

	body := rec.Body.String()
	if !strings.Contains(body, `class="breadcrumb"`) {
		t.Error("README auto-serve should have breadcrumb element")
	}
	if !strings.Contains(body, `<a href="/">root</a>`) {
		t.Error("breadcrumb should have root link")
	}
}

func TestTrailingSlashDirectoryBreadcrumb(t *testing.T) {
	h := handler.New(setupTestDir(t))
	rec := request(h, "/sub/")

	body := rec.Body.String()
	if !strings.Contains(body, `<a href="/">root</a>`) {
		t.Error("trailing-slash directory breadcrumb should have root link")
	}
}

func TestDefaultFilteringHidesNonMarkdown(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "readme.md", "# Hi")
	writeFile(t, dir, "script.py", "print('hello')")
	writeFile(t, dir, "data.json", "{}")

	h := handler.New(dir)
	rec := request(h, "/")
	body := rec.Body.String()

	if !strings.Contains(body, "readme.md") {
		t.Error("expected readme.md to be visible")
	}
	if strings.Contains(body, "script.py") {
		t.Error("expected script.py to be hidden by default")
	}
	if strings.Contains(body, "data.json") {
		t.Error("expected data.json to be hidden by default")
	}
}

func TestDefaultFilteringHidesEmptyDirs(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "doc.md", "# Doc")
	os.MkdirAll(filepath.Join(dir, "empty"), 0o755)
	os.MkdirAll(filepath.Join(dir, "has-md"), 0o755)
	writeFile(t, dir, "has-md/page.md", "# Page")
	os.MkdirAll(filepath.Join(dir, "no-md"), 0o755)
	writeFile(t, dir, "no-md/script.py", "x")

	h := handler.New(dir)
	rec := request(h, "/")
	body := rec.Body.String()

	if !strings.Contains(body, "has-md") {
		t.Error("expected has-md directory to be visible")
	}
	if strings.Contains(body, "empty") {
		t.Error("expected empty directory to be hidden")
	}
	if strings.Contains(body, "no-md") {
		t.Error("expected no-md directory to be hidden")
	}
}

func TestShowAllShowsEverything(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "readme.md", "# Hi")
	writeFile(t, dir, "script.py", "print('hello')")
	os.MkdirAll(filepath.Join(dir, "empty"), 0o755)

	h := handler.New(dir)
	h.ShowAll = true
	rec := request(h, "/")
	body := rec.Body.String()

	if !strings.Contains(body, "readme.md") {
		t.Error("expected readme.md to be visible")
	}
	if !strings.Contains(body, "script.py") {
		t.Error("expected script.py to be visible with ShowAll")
	}
	if !strings.Contains(body, "empty") {
		t.Error("expected empty directory to be visible with ShowAll")
	}
}

func TestAutoReadmeStillWorksWithFiltering(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "docs"), 0o755)
	writeFile(t, dir, "docs/README.md", "# Auto Readme")

	h := handler.New(dir)
	h.AutoReadme = true
	rec := request(h, "/docs")

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Auto Readme") {
		t.Error("expected README.md content with AutoReadme enabled")
	}
}
