package cli_test

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/driangle/viewmd/apps/cli/internal/testutil"
)

// ---------------------------------------------------------------------------
// Text file serving
// ---------------------------------------------------------------------------

func TestTextFileWrappedInPre(t *testing.T) {
	srv := testutil.StartServer(t, map[string]string{
		"hello.py": "print('hello')\n",
	})
	defer srv.Close()

	body := testutil.Get(t, srv.URL, "/hello.py")
	if !strings.Contains(body, "<pre>") {
		t.Error("expected <pre> tag in text file response")
	}
}

func TestTextFileShowsFilename(t *testing.T) {
	srv := testutil.StartServer(t, map[string]string{
		"hello.py": "print('hello')\n",
	})
	defer srv.Close()

	body := testutil.Get(t, srv.URL, "/hello.py")
	if !strings.Contains(body, "hello.py") {
		t.Error("expected filename in text file response")
	}
}

func TestTextFileContentEscaped(t *testing.T) {
	srv := testutil.StartServer(t, map[string]string{
		"hello.py": "print('hello')\n",
	})
	defer srv.Close()

	body := testutil.Get(t, srv.URL, "/hello.py")
	if !strings.Contains(body, "print(&#39;hello&#39;)") {
		t.Error("expected HTML-escaped single quotes in text file response")
	}
}

func TestDotfileServedAsText(t *testing.T) {
	srv := testutil.StartServer(t, map[string]string{
		".gitignore": "__pycache__/\n",
	})
	defer srv.Close()

	body := testutil.Get(t, srv.URL, "/.gitignore")
	if !strings.Contains(body, "<pre>") {
		t.Error("expected <pre> tag for dotfile")
	}
	if !strings.Contains(body, "__pycache__/") {
		t.Error("expected dotfile content in response")
	}
}

func TestKnownFilenameServedAsText(t *testing.T) {
	srv := testutil.StartServer(t, map[string]string{
		"Makefile": "all:\n\techo hi\n",
	})
	defer srv.Close()

	body := testutil.Get(t, srv.URL, "/Makefile")
	if !strings.Contains(body, "<pre>") {
		t.Error("expected <pre> tag for Makefile")
	}
}

func TestBinaryWithTextExtDoesNotCrash(t *testing.T) {
	srv := testutil.StartServer(t, nil)
	defer srv.Close()

	// Write invalid UTF-8 bytes directly (can't use the string-based extra map)
	err := os.WriteFile(
		filepath.Join(srv.TempDir, "broken.txt"),
		[]byte{0x80, 0x81, 0x82, 0xff},
		0o644,
	)
	if err != nil {
		t.Fatalf("writing broken.txt: %v", err)
	}

	resp := testutil.GetResponse(t, srv.URL, "/broken.txt")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 for broken.txt, got %d", resp.StatusCode)
	}
}

// ---------------------------------------------------------------------------
// Directory listing
// ---------------------------------------------------------------------------

func TestRootListsMarkdownAndDirsWithMarkdown(t *testing.T) {
	srv := testutil.StartServer(t, map[string]string{
		"hello.py":       "print('hello')\n",
		"docs/README.md": "# Docs\n",
	})
	defer srv.Close()

	// Create empty_dir (no content needed)
	os.Mkdir(filepath.Join(srv.TempDir, "empty_dir"), 0o755)

	body := testutil.Get(t, srv.URL, "/")

	// Markdown files should appear
	if !strings.Contains(body, "basic.md") {
		t.Error("expected basic.md in root listing")
	}
	// Dirs with markdown should appear
	if !strings.Contains(body, "docs") {
		t.Error("expected docs/ subdirectory in root listing")
	}
	// Non-markdown files hidden by default
	if strings.Contains(body, "hello.py") {
		t.Error("expected hello.py to be hidden by default filtering")
	}
	// Empty dirs hidden by default
	if strings.Contains(body, "empty_dir") {
		t.Error("expected empty_dir to be hidden by default filtering")
	}
}

func TestRootListsAllWithShowAll(t *testing.T) {
	srv := testutil.StartServerWithShowAll(t, map[string]string{
		"hello.py":       "print('hello')\n",
		"docs/README.md": "# Docs\n",
	})
	defer srv.Close()

	os.Mkdir(filepath.Join(srv.TempDir, "empty_dir"), 0o755)

	body := testutil.Get(t, srv.URL, "/")
	if !strings.Contains(body, "hello.py") {
		t.Error("expected hello.py in root listing with ShowAll")
	}
	if !strings.Contains(body, "basic.md") {
		t.Error("expected basic.md in root listing with ShowAll")
	}
	if !strings.Contains(body, "docs") {
		t.Error("expected docs/ in root listing with ShowAll")
	}
}

func TestSubdirWithReadmeShowsListingByDefault(t *testing.T) {
	srv := testutil.StartServer(t, map[string]string{
		"docs/README.md": "# Docs\n",
	})
	defer srv.Close()

	body := testutil.Get(t, srv.URL, "/docs")
	if !strings.Contains(body, "Directory listing") {
		t.Error("expected directory listing by default, not README content")
	}
	if !strings.Contains(body, "README.md") {
		t.Error("expected README.md in directory listing")
	}
}

func TestSubdirWithoutReadmeShowsListing(t *testing.T) {
	srv := testutil.StartServer(t, nil)
	defer srv.Close()

	os.Mkdir(filepath.Join(srv.TempDir, "empty_dir"), 0o755)

	body := testutil.Get(t, srv.URL, "/empty_dir")
	if !strings.Contains(body, "Directory") {
		t.Error("expected 'Directory' in listing for dir without README")
	}
}

func TestSubdirListingHasBreadcrumb(t *testing.T) {
	srv := testutil.StartServer(t, nil)
	defer srv.Close()

	os.Mkdir(filepath.Join(srv.TempDir, "empty_dir"), 0o755)

	body := testutil.Get(t, srv.URL, "/empty_dir")
	if !strings.Contains(body, `class="breadcrumb"`) {
		t.Error("expected breadcrumb in subdirectory listing")
	}
	if !strings.Contains(body, `<a href="/">root</a>`) {
		t.Error("expected root link in breadcrumb")
	}
}

// ---------------------------------------------------------------------------
// Routing
// ---------------------------------------------------------------------------

func TestMarkdownRoutedToRenderer(t *testing.T) {
	srv := testutil.StartServer(t, nil)
	defer srv.Close()

	body := testutil.Get(t, srv.URL, "/plain.md")
	if !strings.Contains(body, "<h1>Just a heading</h1>") {
		t.Error("expected rendered heading for markdown file")
	}
}

func TestBinaryFileServedRaw(t *testing.T) {
	pngHeader := append([]byte("\x89PNG\r\n\x1a\n"), make([]byte, 20)...)
	srv := testutil.StartServer(t, nil)
	defer srv.Close()

	os.WriteFile(filepath.Join(srv.TempDir, "image.png"), pngHeader, 0o644)

	resp := testutil.GetResponse(t, srv.URL, "/image.png")
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading image.png response: %v", err)
	}
	if len(data) < 4 || string(data[:4]) != "\x89PNG" {
		t.Error("expected raw PNG magic bytes in response")
	}
}

func TestMissingFileReturns404(t *testing.T) {
	srv := testutil.StartServer(t, nil)
	defer srv.Close()

	resp := testutil.GetResponse(t, srv.URL, "/nonexistent.txt")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 for missing file, got %d", resp.StatusCode)
	}
}

func TestURLEncodedPathWorks(t *testing.T) {
	srv := testutil.StartServer(t, map[string]string{
		"hello.py": "print('hello')\n",
	})
	defer srv.Close()

	body := testutil.Get(t, srv.URL, "/hello.py")
	if !strings.Contains(body, "<pre>") {
		t.Error("expected <pre> tag for URL-encoded path")
	}
}
