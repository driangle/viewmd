package cli_test

import (
	"html"
	"strings"
	"testing"

	"github.com/driangle/viewmd/apps/cli/internal/testutil"
)

func TestFrontmatterRenderedAsTable(t *testing.T) {
	srv := testutil.StartServer(t, nil)
	defer srv.Close()

	body := testutil.Get(t, srv.URL, "/basic.md")

	checks := []string{
		`<div class="frontmatter">`,
		"<table>",
		`class="fm-key"`,
		"title",
		"Hello",
		"2024-01-01",
		"<h1>Heading</h1>",
	}
	for _, want := range checks {
		if !strings.Contains(body, want) {
			t.Errorf("expected body to contain %q", want)
		}
	}
}

func TestNoFrontmatterNoTable(t *testing.T) {
	srv := testutil.StartServer(t, nil)
	defer srv.Close()

	body := testutil.Get(t, srv.URL, "/plain.md")

	if strings.Contains(body, `<div class="frontmatter">`) {
		t.Error("expected no frontmatter div for plain.md")
	}
	if !strings.Contains(body, "<h1>Just a heading</h1>") {
		t.Error("expected rendered heading in plain.md")
	}
}

func TestFrontmatterHTMLEscaping(t *testing.T) {
	srv := testutil.StartServer(t, nil)
	defer srv.Close()

	body := testutil.Get(t, srv.URL, "/escaping.md")

	xssPayload := `<script>alert("xss")</script>`
	if strings.Contains(body, xssPayload) {
		t.Error("expected frontmatter XSS payload to be escaped")
	}

	escaped := html.EscapeString(xssPayload)
	if !strings.Contains(body, escaped) {
		t.Errorf("expected escaped script tag %q in body", escaped)
	}
}

func TestFrontmatterStrippedFromBody(t *testing.T) {
	srv := testutil.StartServer(t, nil)
	defer srv.Close()

	body := testutil.Get(t, srv.URL, "/stripped.md")

	// Extract the rendered body between the frontmatter div and the raw-content textarea.
	// The raw textarea intentionally contains the full file (including frontmatter) for copying.
	parts := strings.SplitN(body, "</div>", 2)
	afterFM := parts[len(parts)-1]
	if idx := strings.Index(afterFM, `<textarea id="raw-content"`); idx >= 0 {
		afterFM = afterFM[:idx]
	}

	if strings.Contains(afterFM, "status: draft") {
		t.Error("frontmatter YAML should not appear in the rendered body")
	}
}

func TestFrontmatterCSSPresent(t *testing.T) {
	srv := testutil.StartServer(t, nil)
	defer srv.Close()

	body := testutil.Get(t, srv.URL, "/basic.md")

	if !strings.Contains(body, ".frontmatter") {
		t.Error("expected .frontmatter CSS class in output")
	}
	if !strings.Contains(body, ".fm-key") {
		t.Error("expected .fm-key CSS class in output")
	}
}
