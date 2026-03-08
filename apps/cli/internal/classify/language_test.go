package classify

import "testing"

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"main.go", "go"},
		{"script.py", "python"},
		{"app.js", "javascript"},
		{"index.tsx", "typescript"},
		{"style.css", "css"},
		{"config.yaml", "yaml"},
		{"data.json", "json"},
		{"query.sql", "sql"},
		{"run.sh", "bash"},
		{"Makefile", "makefile"},
		{"Dockerfile", "dockerfile"},
		{"Gemfile", "ruby"},
		{"notes.txt", ""},
		{"output.log", ""},
		{"unknown.xyz", ""},
		{".bashrc", ""},
		{"MAIN.GO", "go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectLanguage(tt.name)
			if got != tt.want {
				t.Errorf("DetectLanguage(%q) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}
