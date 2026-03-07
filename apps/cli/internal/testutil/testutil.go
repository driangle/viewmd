// Package testutil provides shared helpers for integration tests.
package testutil

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/driangle/viewmd/apps/cli/internal/handler"
)

// TestdataDir returns the absolute path to the testdata/ directory.
func TestdataDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata")
}

// Server holds a running test server's details and cleanup function.
type Server struct {
	URL     string
	TempDir string
	close   func()
}

// Close shuts down the server and removes the temp directory.
func (s *Server) Close() {
	s.close()
}

// StartServer copies fixtures from testdata/ plus any extra files into a temp
// directory, starts a real HTTP server on a random port, and returns a Server.
// Extra files are specified as path->content pairs relative to the temp dir.
func StartServer(t *testing.T, extra map[string]string) *Server {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "viewmd-test-*")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}

	copyFixtures(t, tmpDir)

	for relPath, content := range extra {
		full := filepath.Join(tmpDir, relPath)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("creating dir for %s: %v", relPath, err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("writing extra file %s: %v", relPath, err)
		}
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listening on random port: %v", err)
	}

	srv := &http.Server{Handler: handler.New(tmpDir)}
	go srv.Serve(listener)

	baseURL := fmt.Sprintf("http://127.0.0.1:%d", listener.Addr().(*net.TCPAddr).Port)

	return &Server{
		URL:     baseURL,
		TempDir: tmpDir,
		close: func() {
			srv.Close()
			os.RemoveAll(tmpDir)
		},
	}
}

// Get fetches baseURL+path and returns the response body as a string.
func Get(t *testing.T, baseURL, path string) string {
	t.Helper()

	resp, err := http.Get(baseURL + path)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading response body for %s: %v", path, err)
	}
	return string(body)
}

// GetResponse fetches baseURL+path and returns the full *http.Response.
// The caller is responsible for closing the body.
func GetResponse(t *testing.T, baseURL, path string) *http.Response {
	t.Helper()

	resp, err := http.Get(baseURL + path)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	return resp
}

func copyFixtures(t *testing.T, destDir string) {
	t.Helper()

	srcDir := TestdataDir()
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		t.Fatalf("reading testdata dir: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		src := filepath.Join(srcDir, entry.Name())
		dst := filepath.Join(destDir, entry.Name())

		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("reading fixture %s: %v", entry.Name(), err)
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			t.Fatalf("writing fixture %s: %v", entry.Name(), err)
		}
	}
}
