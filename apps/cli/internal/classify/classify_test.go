package classify_test

import (
	"testing"

	"github.com/driangle/viewmd/apps/cli/internal/classify"
)

func TestIsTextFile(t *testing.T) {
	tests := []struct {
		name string
		file string
		want bool
	}{
		{"python source", "script.py", true},
		{"json data", "data.json", true},
		{"png image", "photo.png", false},
		{"pdf document", "doc.pdf", false},
		{"makefile", "Makefile", true},
		{"license uppercase", "LICENSE", true},
		{"dotfile with extension", ".env", true},
		{"gitignore", ".gitignore", true},
		{"unknown file no extension", "randomfile", false},
		{"markdown not text", "readme.md", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classify.IsTextFile(tt.file)
			if got != tt.want {
				t.Errorf("IsTextFile(%q) = %v, want %v", tt.file, got, tt.want)
			}
		})
	}
}
