"""Tests for frontmatter parsing and rendering in viewmd."""
import html
import os
import shutil
from http.server import HTTPServer
from pathlib import Path
from threading import Thread
from urllib.request import urlopen

import pytest

from viewmd import MarkdownHandler, parse_frontmatter

FIXTURES = Path(__file__).parent / "fixtures"


# ---------------------------------------------------------------------------
# parse_frontmatter unit tests
# ---------------------------------------------------------------------------

class TestParseFrontmatter:
    def test_basic(self):
        fm, body = parse_frontmatter("---\ntitle: Hello\n---\n# Body")
        assert fm == {"title": "Hello"}
        assert body == "\n# Body"

    def test_multiple_keys(self):
        content = "---\ntitle: Hello\ndate: 2024-01-01\ntags: a, b\n---\nBody"
        fm, body = parse_frontmatter(content)
        assert fm == {"title": "Hello", "date": "2024-01-01", "tags": "a, b"}
        assert body == "\nBody"

    def test_no_frontmatter(self):
        content = "# Just markdown\nSome text"
        fm, body = parse_frontmatter(content)
        assert fm is None
        assert body == content

    def test_empty_string(self):
        fm, body = parse_frontmatter("")
        assert fm is None
        assert body == ""

    def test_single_delimiter(self):
        content = "---\nbroken"
        fm, body = parse_frontmatter(content)
        assert fm is None
        assert body == content

    def test_value_with_colons(self):
        content = "---\nurl: https://example.com\ntime: 10:30:00\n---\nBody"
        fm, body = parse_frontmatter(content)
        assert fm == {"url": "https://example.com", "time": "10:30:00"}

    def test_whitespace_around_keys_and_values(self):
        content = "---\n  title  :  Hello World  \n---\nBody"
        fm, _ = parse_frontmatter(content)
        assert fm == {"title": "Hello World"}

    def test_empty_value(self):
        content = "---\ndraft:\n---\nBody"
        fm, _ = parse_frontmatter(content)
        assert fm == {"draft": ""}

    def test_lines_without_colon_skipped(self):
        content = "---\ntitle: Hello\njust a line\ndate: 2024\n---\nBody"
        fm, _ = parse_frontmatter(content)
        assert fm == {"title": "Hello", "date": "2024"}
        assert "just a line" not in fm

    def test_body_preserved_after_frontmatter(self):
        body_text = "\n# Heading\n\nParagraph with `code` and **bold**."
        content = f"---\ntitle: T\n---{body_text}"
        _, body = parse_frontmatter(content)
        assert body == body_text

    def test_triple_dash_in_body(self):
        """A --- in the body (after frontmatter) should be left in the body."""
        content = "---\ntitle: T\n---\nBody\n---\nMore body"
        fm, body = parse_frontmatter(content)
        assert fm == {"title": "T"}
        assert "---" in body
        assert "More body" in body

    def test_only_delimiters(self):
        content = "---\n---\nBody"
        fm, body = parse_frontmatter(content)
        assert fm == {}
        assert body == "\nBody"

    def test_content_starting_with_dashes_but_not_frontmatter(self):
        content = "---- not frontmatter\nstuff"
        fm, body = parse_frontmatter(content)
        assert fm is None
        assert body == content


# ---------------------------------------------------------------------------
# serve_markdown integration tests (via a real HTTP server)
# ---------------------------------------------------------------------------

@pytest.fixture()
def server(tmp_path):
    """Copy fixtures into a temp dir and start a viewmd server from there."""
    for f in FIXTURES.iterdir():
        shutil.copy(f, tmp_path / f.name)

    old_cwd = os.getcwd()
    os.chdir(tmp_path)

    srv = HTTPServer(("127.0.0.1", 0), MarkdownHandler)
    port = srv.server_address[1]
    thread = Thread(target=srv.serve_forever, daemon=True)
    thread.start()
    yield f"http://127.0.0.1:{port}"
    srv.shutdown()
    os.chdir(old_cwd)


def _get(base_url, path):
    """Fetch a URL path and return the decoded response body."""
    with urlopen(f"{base_url}/{path}") as resp:
        return resp.read().decode("utf-8")


class TestServeMarkdownFrontmatter:
    def test_frontmatter_rendered_as_table(self, server):
        body = _get(server, "basic.md")
        assert '<div class="frontmatter">' in body
        assert "<table>" in body
        assert 'class="fm-key"' in body
        assert "title" in body
        assert "Hello" in body
        assert "2024-01-01" in body
        assert "<h1>Heading</h1>" in body

    def test_no_frontmatter_no_table(self, server):
        body = _get(server, "plain.md")
        assert '<div class="frontmatter">' not in body
        assert "<h1>Just a heading</h1>" in body

    def test_frontmatter_html_escaping(self, server):
        body = _get(server, "escaping.md")
        assert "<script>" not in body
        assert html.escape('<script>alert("xss")</script>') in body

    def test_frontmatter_stripped_from_body(self, server):
        """Frontmatter YAML lines should not leak into the rendered markdown body."""
        body = _get(server, "stripped.md")
        after_fm = body.split("</div>", 1)[-1]
        assert "status: draft" not in after_fm

    def test_frontmatter_css_present(self, server):
        body = _get(server, "basic.md")
        assert ".frontmatter" in body
        assert ".fm-key" in body
