package handler

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsIgnored(t *testing.T) {
	tests := []struct {
		name     string
		relPath  string
		patterns []string
		want     bool
	}{
		{".git", ".git", []string{".git"}, true},
		{"node_modules", "node_modules", []string{"node_modules"}, true},
		{"README.md", "README.md", []string{".git", "node_modules"}, false},
		{"build.log", "build.log", []string{"*.log"}, true},
		{"main.go", "main.go", []string{"*.log"}, false},
		{".git", ".git", nil, false},
		{".git", ".git", []string{}, false},
		{".gitignore", ".gitignore", []string{".git"}, false},
		// Path-based patterns
		{"worktrees", ".claude/worktrees", []string{".claude/worktrees"}, true},
		{"foo.txt", ".claude/worktrees/foo.txt", []string{".claude/worktrees"}, true},
		{"worktrees", "other/worktrees", []string{".claude/worktrees"}, false},
		{"release.yml", "docs/release.yml", []string{"docs/*.yml"}, true},
		{"release.yml", "other/release.yml", []string{"docs/*.yml"}, false},
		// Doublestar patterns
		{"test.log", "a/b/test.log", []string{"**/*.log"}, true},
		{"test.log", "test.log", []string{"**/*.log"}, true},
		{"test.go", "a/b/test.go", []string{"**/*.log"}, false},
	}
	for _, tt := range tests {
		if got := isIgnored(tt.name, tt.relPath, tt.patterns); got != tt.want {
			t.Errorf("isIgnored(%q, %q, %v) = %v, want %v", tt.name, tt.relPath, tt.patterns, got, tt.want)
		}
	}
}

func TestIsMarkdownFile(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"README.md", true},
		{"notes.markdown", true},
		{"doc.MD", true},
		{"file.MARKDOWN", true},
		{"script.py", false},
		{"image.png", false},
		{".gitignore", false},
		{"noext", false},
	}
	for _, tt := range tests {
		if got := isMarkdownFile(tt.name); got != tt.want {
			t.Errorf("isMarkdownFile(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestHasMarkdownFiles(t *testing.T) {
	tests := []struct {
		name  string
		setup func(dir string)
		want  bool
	}{
		{
			name: "directory with markdown file",
			setup: func(dir string) {
				os.WriteFile(filepath.Join(dir, "doc.md"), []byte("# Hi"), 0o644)
			},
			want: true,
		},
		{
			name: "directory with only non-markdown files",
			setup: func(dir string) {
				os.WriteFile(filepath.Join(dir, "script.py"), []byte("x"), 0o644)
			},
			want: false,
		},
		{
			name: "empty directory",
			setup: func(dir string) {},
			want:  false,
		},
		{
			name: "nested markdown file",
			setup: func(dir string) {
				sub := filepath.Join(dir, "sub")
				os.MkdirAll(sub, 0o755)
				os.WriteFile(filepath.Join(sub, "page.md"), []byte("# Page"), 0o644)
			},
			want: true,
		},
		{
			name: "deeply nested markdown file",
			setup: func(dir string) {
				deep := filepath.Join(dir, "a", "b", "c")
				os.MkdirAll(deep, 0o755)
				os.WriteFile(filepath.Join(deep, "deep.markdown"), []byte("# Deep"), 0o644)
			},
			want: true,
		},
		{
			name: "nested dirs with no markdown",
			setup: func(dir string) {
				sub := filepath.Join(dir, "sub")
				os.MkdirAll(sub, 0o755)
				os.WriteFile(filepath.Join(sub, "data.json"), []byte("{}"), 0o644)
			},
			want: false,
		},
		{
			name:  "nonexistent directory",
			setup: func(dir string) {},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.name == "nonexistent directory" {
				dir = filepath.Join(dir, "nope")
			} else {
				tt.setup(dir)
			}
			if got := hasMarkdownFiles(dir); got != tt.want {
				t.Errorf("hasMarkdownFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
