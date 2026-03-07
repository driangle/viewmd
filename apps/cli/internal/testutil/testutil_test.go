package testutil

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestStartServerCopiesFixtures(t *testing.T) {
	srv := StartServer(t, nil)
	defer srv.Close()

	for _, name := range []string{"basic.md", "escaping.md", "plain.md", "stripped.md"} {
		path := filepath.Join(srv.TempDir, name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected fixture %s in temp dir, got error: %v", name, err)
		}
	}
}

func TestStartServerWithExtraFiles(t *testing.T) {
	srv := StartServer(t, map[string]string{
		"hello.txt":          "world",
		"subdir/nested.txt":  "nested content",
	})
	defer srv.Close()

	data, err := os.ReadFile(filepath.Join(srv.TempDir, "hello.txt"))
	if err != nil {
		t.Fatalf("reading extra file: %v", err)
	}
	if string(data) != "world" {
		t.Errorf("expected %q, got %q", "world", string(data))
	}

	data, err = os.ReadFile(filepath.Join(srv.TempDir, "subdir", "nested.txt"))
	if err != nil {
		t.Fatalf("reading nested extra file: %v", err)
	}
	if string(data) != "nested content" {
		t.Errorf("expected %q, got %q", "nested content", string(data))
	}
}

func TestServerResponds(t *testing.T) {
	srv := StartServer(t, nil)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/basic.md")
	if err != nil {
		t.Fatalf("GET /basic.md: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestGetReturnsBody(t *testing.T) {
	srv := StartServer(t, nil)
	defer srv.Close()

	body := Get(t, srv.URL, "/plain.md")
	if body == "" {
		t.Error("expected non-empty response body")
	}
}
