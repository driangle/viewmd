package handler_test

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/driangle/viewmd/apps/cli/internal/handler"
)

type searchResponse struct {
	Results []struct {
		Path    string `json:"path"`
		Name    string `json:"name"`
		IsDir   bool   `json:"isDir"`
		Match   string `json:"match"`
		Snippet string `json:"snippet"`
	} `json:"results"`
}

func setupSearchDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	writeFile(t, dir, "hello.md", "# Hello World\nThis is a test document.")
	writeFile(t, dir, "guide.md", "# User Guide\nHow to use the app.")
	writeFile(t, dir, "notes.txt", "Some notes about the project.")
	writeFile(t, dir, "code.go", "package main\nfunc main() {}")

	os.MkdirAll(filepath.Join(dir, "docs"), 0o755)
	writeFile(t, dir, "docs/setup.md", "# Setup\nInstallation steps here.")
	os.MkdirAll(filepath.Join(dir, ".git"), 0o755)
	writeFile(t, dir, ".git/config", "hidden")

	return dir
}

func searchRequest(h http.Handler, query string) searchResponse {
	rec := request(h, "/-/search?"+query)
	var resp searchResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	return resp
}

func TestSearchEmptyQuery(t *testing.T) {
	h := handler.New(setupSearchDir(t))
	h.ShowAll = true
	resp := searchRequest(h, "q=")

	if len(resp.Results) != 0 {
		t.Errorf("expected 0 results for empty query, got %d", len(resp.Results))
	}
}

func TestSearchByName(t *testing.T) {
	h := handler.New(setupSearchDir(t))
	h.ShowAll = true
	resp := searchRequest(h, "q=hello&mode=name")

	if len(resp.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(resp.Results))
	}
	if resp.Results[0].Name != "hello.md" {
		t.Errorf("expected hello.md, got %s", resp.Results[0].Name)
	}
	if resp.Results[0].Match != "name" {
		t.Errorf("expected match type 'name', got %s", resp.Results[0].Match)
	}
}

func TestSearchByContent(t *testing.T) {
	h := handler.New(setupSearchDir(t))
	h.ShowAll = true
	resp := searchRequest(h, "q=installation&mode=content")

	if len(resp.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(resp.Results))
	}
	if resp.Results[0].Path != "docs/setup.md" {
		t.Errorf("expected docs/setup.md, got %s", resp.Results[0].Path)
	}
	if resp.Results[0].Match != "content" {
		t.Errorf("expected match type 'content', got %s", resp.Results[0].Match)
	}
	if resp.Results[0].Snippet == "" {
		t.Error("expected a snippet for content match")
	}
}

func TestSearchBothMode(t *testing.T) {
	h := handler.New(setupSearchDir(t))
	h.ShowAll = true
	// "guide" matches filename and "guide" appears in content of guide.md
	resp := searchRequest(h, "q=guide&mode=both")

	hasBoth := false
	for _, r := range resp.Results {
		if r.Name == "guide.md" && r.Match == "both" {
			hasBoth = true
		}
	}
	if !hasBoth {
		t.Error("expected guide.md to match as 'both' (name + content)")
	}
	// Should not have duplicate entries for the same file
	paths := map[string]int{}
	for _, r := range resp.Results {
		paths[r.Path]++
	}
	for p, count := range paths {
		if count > 1 {
			t.Errorf("duplicate result for %s: appeared %d times", p, count)
		}
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	h := handler.New(setupSearchDir(t))
	h.ShowAll = true
	resp := searchRequest(h, "q=HELLO&mode=name")

	if len(resp.Results) != 1 {
		t.Fatalf("expected 1 result for case-insensitive search, got %d", len(resp.Results))
	}
}

func TestSearchSkipsHiddenDirs(t *testing.T) {
	h := handler.New(setupSearchDir(t))
	h.ShowAll = true
	resp := searchRequest(h, "q=config&mode=both")

	for _, r := range resp.Results {
		if r.Name == "config" {
			t.Error("should not return files from .git directory")
		}
	}
}

func TestSearchRespectsShowAll(t *testing.T) {
	h := handler.New(setupSearchDir(t))
	// ShowAll=false (default): only markdown files and dirs
	resp := searchRequest(h, "q=notes&mode=name")

	for _, r := range resp.Results {
		if r.Name == "notes.txt" {
			t.Error("should not return .txt files when ShowAll is false")
		}
	}
}

func TestSearchReturnsJSON(t *testing.T) {
	h := handler.New(setupSearchDir(t))
	rec := request(h, "/-/search?q=hello")

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestSearchResultLimit(t *testing.T) {
	dir := t.TempDir()
	for i := range 60 {
		writeFile(t, dir, filepath.Join("match_"+string(rune('a'+i/26))+string(rune('a'+i%26))+".md"), "# Match")
	}

	h := handler.New(dir)
	resp := searchRequest(h, "q=match&mode=name")

	if len(resp.Results) > 50 {
		t.Errorf("expected at most 50 results, got %d", len(resp.Results))
	}
}

func TestSearchDirMatch(t *testing.T) {
	h := handler.New(setupSearchDir(t))
	h.ShowAll = true
	resp := searchRequest(h, "q=docs&mode=name")

	found := false
	for _, r := range resp.Results {
		if r.Name == "docs" && r.IsDir {
			found = true
		}
	}
	if !found {
		t.Error("expected to find 'docs' directory in results")
	}
}
