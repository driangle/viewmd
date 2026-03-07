package frontmatter_test

import (
	"testing"

	"github.com/driangle/viewmd/apps/cli/internal/frontmatter"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantMeta map[string]string
		wantBody string
	}{
		{
			name:     "basic",
			input:    "---\ntitle: Hello\n---\n# Body",
			wantMeta: map[string]string{"title": "Hello"},
			wantBody: "\n# Body",
		},
		{
			name:  "multiple keys",
			input: "---\ntitle: Hello\ndate: 2024-01-01\ntags: a, b\n---\nBody",
			wantMeta: map[string]string{
				"title": "Hello",
				"date":  "2024-01-01",
				"tags":  "a, b",
			},
			wantBody: "\nBody",
		},
		{
			name:     "no frontmatter",
			input:    "# Just markdown\nSome text",
			wantMeta: nil,
			wantBody: "# Just markdown\nSome text",
		},
		{
			name:     "empty string",
			input:    "",
			wantMeta: nil,
			wantBody: "",
		},
		{
			name:     "single delimiter",
			input:    "---\nbroken",
			wantMeta: nil,
			wantBody: "---\nbroken",
		},
		{
			name:  "value with colons",
			input: "---\nurl: https://example.com\ntime: 10:30:00\n---\nBody",
			wantMeta: map[string]string{
				"url":  "https://example.com",
				"time": "10:30:00",
			},
			wantBody: "\nBody",
		},
		{
			name:     "whitespace around keys and values",
			input:    "---\n  title  :  Hello World  \n---\nBody",
			wantMeta: map[string]string{"title": "Hello World"},
			wantBody: "\nBody",
		},
		{
			name:     "empty value",
			input:    "---\ndraft:\n---\nBody",
			wantMeta: map[string]string{"draft": ""},
			wantBody: "\nBody",
		},
		{
			name:  "lines without colon skipped",
			input: "---\ntitle: Hello\njust a line\ndate: 2024\n---\nBody",
			wantMeta: map[string]string{
				"title": "Hello",
				"date":  "2024",
			},
			wantBody: "\nBody",
		},
		{
			name:     "body preserved",
			input:    "---\ntitle: T\n---\nLine 1\nLine 2\n\nLine 4",
			wantMeta: map[string]string{"title": "T"},
			wantBody: "\nLine 1\nLine 2\n\nLine 4",
		},
		{
			name:     "triple dash in body",
			input:    "---\ntitle: T\n---\nBody\n---\nMore body",
			wantMeta: map[string]string{"title": "T"},
			wantBody: "\nBody\n---\nMore body",
		},
		{
			name:     "only delimiters",
			input:    "---\n---\nBody",
			wantMeta: map[string]string{},
			wantBody: "\nBody",
		},
		{
			name:     "content starting with dashes but not frontmatter",
			input:    "---- not frontmatter\nstuff",
			wantMeta: nil,
			wantBody: "---- not frontmatter\nstuff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta, body := frontmatter.Parse(tt.input)

			if tt.wantMeta == nil {
				if meta != nil {
					t.Errorf("expected nil meta, got %v", meta)
				}
			} else {
				if meta == nil {
					t.Fatalf("expected meta %v, got nil", tt.wantMeta)
				}
				if len(meta) != len(tt.wantMeta) {
					t.Errorf("meta length = %d, want %d\n  got:  %v\n  want: %v",
						len(meta), len(tt.wantMeta), meta, tt.wantMeta)
				}
				for k, want := range tt.wantMeta {
					if got, ok := meta[k]; !ok {
						t.Errorf("missing key %q", k)
					} else if got != want {
						t.Errorf("meta[%q] = %q, want %q", k, got, want)
					}
				}
			}

			if body != tt.wantBody {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
		})
	}
}
