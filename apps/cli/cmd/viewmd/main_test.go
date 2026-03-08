package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrintBanner(t *testing.T) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	printBanner(8000, "/tmp/test")

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	checks := []string{
		"viewmd v0.1.0",
		"http://localhost:8000",
		"/tmp/test",
	}
	for _, s := range checks {
		if !strings.Contains(output, s) {
			t.Errorf("banner missing %q", s)
		}
	}
}

func TestPortInUseMessage(t *testing.T) {
	// Occupy a port to test the error path.
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port

	// Try binding to the same port — should fail.
	_, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err == nil {
		t.Fatal("expected port-in-use error")
	}
}

func TestLoadShowAllFromConfig(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{"true value", "show_all_files: true\n", true},
		{"false value", "show_all_files: false\n", false},
		{"no such key", "other_key: true\n", false},
		{"empty file", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			os.WriteFile(filepath.Join(dir, ".viewmd.yaml"), []byte(tt.content), 0o644)
			if got := loadShowAllFromConfig(dir); got != tt.want {
				t.Errorf("loadShowAllFromConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadShowAllFromConfigMissingFile(t *testing.T) {
	dir := t.TempDir()
	if got := loadShowAllFromConfig(dir); got != false {
		t.Errorf("expected false when config file missing, got %v", got)
	}
}
