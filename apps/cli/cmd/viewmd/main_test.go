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

func TestBuildIgnorePatternsDefault(t *testing.T) {
	dir := t.TempDir()
	got := buildIgnorePatterns("", dir)
	if len(got) != 1 || got[0] != ".git" {
		t.Errorf("expected default [.git], got %v", got)
	}
}

func TestBuildIgnorePatternsFromFlag(t *testing.T) {
	dir := t.TempDir()
	got := buildIgnorePatterns("node_modules, *.log", dir)
	if len(got) != 2 || got[0] != "node_modules" || got[1] != "*.log" {
		t.Errorf("expected [node_modules *.log], got %v", got)
	}
}

func TestBuildIgnorePatternsFromYAML(t *testing.T) {
	dir := t.TempDir()
	yaml := "ignore:\n  - node_modules\n  - .git\n  - \"*.log\"\n"
	os.WriteFile(filepath.Join(dir, ".viewmd.yaml"), []byte(yaml), 0o644)

	got := buildIgnorePatterns("", dir)
	if len(got) != 3 {
		t.Fatalf("expected 3 patterns, got %v", got)
	}
	if got[0] != "node_modules" || got[1] != ".git" || got[2] != "*.log" {
		t.Errorf("unexpected patterns: %v", got)
	}
}

func TestBuildIgnorePatternsMergesFlagAndYAML(t *testing.T) {
	dir := t.TempDir()
	yaml := "ignore:\n  - .git\n"
	os.WriteFile(filepath.Join(dir, ".viewmd.yaml"), []byte(yaml), 0o644)

	got := buildIgnorePatterns("node_modules", dir)
	if len(got) != 2 || got[0] != "node_modules" || got[1] != ".git" {
		t.Errorf("expected [node_modules .git], got %v", got)
	}
}

func TestLoadIgnoreFromConfig(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			"ignore list",
			"ignore:\n  - node_modules\n  - .git\n",
			[]string{"node_modules", ".git"},
		},
		{
			"no ignore key",
			"show_all_files: true\n",
			nil,
		},
		{
			"empty file",
			"",
			nil,
		},
		{
			"ignore with other keys",
			"show_all_files: true\nignore:\n  - dist\nother: val\n",
			[]string{"dist"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			os.WriteFile(filepath.Join(dir, ".viewmd.yaml"), []byte(tt.content), 0o644)
			got := loadIgnoreFromConfig(dir)
			if len(got) != len(tt.want) {
				t.Errorf("loadIgnoreFromConfig() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("loadIgnoreFromConfig()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
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
