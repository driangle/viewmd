package wikilink

import "testing"

func TestResolve(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		baseURL string
		want    string
	}{
		{
			name:    "simple link",
			input:   "See [[page-name]] for details.",
			baseURL: "/docs/",
			want:    "See [page-name](/docs/page-name.md) for details.",
		},
		{
			name:    "link with display text",
			input:   "See [[page-name|the page]] here.",
			baseURL: "/docs/",
			want:    "See [the page](/docs/page-name.md) here.",
		},
		{
			name:    "link with .md extension",
			input:   "See [[readme.md]] here.",
			baseURL: "/",
			want:    "See [readme.md](/readme.md) here.",
		},
		{
			name:    "link with .markdown extension",
			input:   "See [[notes.markdown]] here.",
			baseURL: "/",
			want:    "See [notes.markdown](/notes.markdown) here.",
		},
		{
			name:    "link with path",
			input:   "See [[sub/page]] for more.",
			baseURL: "/",
			want:    "See [sub/page](/sub/page.md) for more.",
		},
		{
			name:    "multiple links",
			input:   "Link [[foo]] and [[bar|Bar Page]].",
			baseURL: "/",
			want:    "Link [foo](/foo.md) and [Bar Page](/bar.md).",
		},
		{
			name:    "no links",
			input:   "No wiki links here.",
			baseURL: "/",
			want:    "No wiki links here.",
		},
		{
			name:    "empty target",
			input:   "Empty [[]] link.",
			baseURL: "/",
			want:    "Empty [[]] link.",
		},
		{
			name:    "baseURL without trailing slash",
			input:   "See [[page]].",
			baseURL: "/docs",
			want:    "See [page](/docs/page.md).",
		},
		{
			name:    "empty baseURL",
			input:   "See [[page]].",
			baseURL: "",
			want:    "See [page](page.md).",
		},
		{
			name:    "inside code fence preserved",
			input:   "Text [[real-link]] end.",
			baseURL: "/",
			want:    "Text [real-link](/real-link.md) end.",
		},
		{
			name:    "spaces in target trimmed",
			input:   "See [[ page-name ]] here.",
			baseURL: "/",
			want:    "See [page-name](/page-name.md) here.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Resolve(tt.input, tt.baseURL)
			if got != tt.want {
				t.Errorf("Resolve(%q, %q)\n got: %q\nwant: %q", tt.input, tt.baseURL, got, tt.want)
			}
		})
	}
}
