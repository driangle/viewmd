"""Tests for text file serving, directory listing, and request routing."""
import os
import shutil
from http.server import HTTPServer
from pathlib import Path
from threading import Thread
from urllib.error import HTTPError
from urllib.request import urlopen

import pytest

from viewmd import MarkdownHandler

FIXTURES = Path(__file__).parent / "fixtures"


# ---------------------------------------------------------------------------
# is_text_file unit tests
# ---------------------------------------------------------------------------

class TestIsTextFile:
    """Unit tests for MarkdownHandler.is_text_file."""

    @pytest.fixture(autouse=True)
    def handler(self):
        self.h = object.__new__(MarkdownHandler)

    def test_known_extension_py(self):
        assert self.h.is_text_file(Path("script.py")) is True

    def test_known_extension_json(self):
        assert self.h.is_text_file(Path("data.json")) is True

    def test_unknown_extension_png(self):
        assert self.h.is_text_file(Path("photo.png")) is False

    def test_unknown_extension_pdf(self):
        assert self.h.is_text_file(Path("doc.pdf")) is False

    def test_known_filename_makefile(self):
        assert self.h.is_text_file(Path("Makefile")) is True

    def test_known_filename_case_insensitive(self):
        assert self.h.is_text_file(Path("LICENSE")) is True

    def test_dotfile_detected(self):
        assert self.h.is_text_file(Path(".env")) is True

    def test_gitignore_in_extensions(self):
        assert self.h.is_text_file(Path(".gitignore")) is True

    def test_no_extension_unknown_name(self):
        assert self.h.is_text_file(Path("randomfile")) is False

    def test_markdown_not_in_text_extensions(self):
        assert self.h.is_text_file(Path("readme.md")) is False


# ---------------------------------------------------------------------------
# Integration test fixture
# ---------------------------------------------------------------------------

@pytest.fixture()
def server(tmp_path):
    """Set up a temp dir with various file types and start a server."""
    for f in FIXTURES.iterdir():
        shutil.copy(f, tmp_path / f.name)

    (tmp_path / "hello.py").write_text("print('hello')\n", encoding="utf-8")
    (tmp_path / ".gitignore").write_text("__pycache__/\n", encoding="utf-8")
    (tmp_path / "Makefile").write_text("all:\n\techo hi\n", encoding="utf-8")
    (tmp_path / "image.png").write_bytes(b"\x89PNG\r\n\x1a\n" + b"\x00" * 20)
    # Binary content with a .txt extension to trigger UnicodeDecodeError fallback
    (tmp_path / "broken.txt").write_bytes(b"\x80\x81\x82\xff")

    subdir = tmp_path / "docs"
    subdir.mkdir()
    (subdir / "README.md").write_text("# Docs\n", encoding="utf-8")

    (tmp_path / "empty_dir").mkdir()

    old_cwd = os.getcwd()
    os.chdir(tmp_path)

    srv = HTTPServer(("127.0.0.1", 0), MarkdownHandler)
    port = srv.server_address[1]
    thread = Thread(target=srv.serve_forever, daemon=True)
    thread.start()
    yield f"http://127.0.0.1:{port}"
    srv.shutdown()
    os.chdir(old_cwd)


def _get(base_url, path=""):
    """Fetch a URL path and return the decoded response body."""
    with urlopen(f"{base_url}/{path}") as resp:
        return resp.read().decode("utf-8")


# ---------------------------------------------------------------------------
# serve_text_file integration tests
# ---------------------------------------------------------------------------

class TestServeTextFile:
    def test_text_file_wrapped_in_pre(self, server):
        body = _get(server, "hello.py")
        assert "<pre>" in body

    def test_text_file_shows_filename(self, server):
        body = _get(server, "hello.py")
        assert "hello.py" in body

    def test_text_file_content_escaped(self, server):
        body = _get(server, "hello.py")
        assert "print(&#x27;hello&#x27;)" in body

    def test_dotfile_served_as_text(self, server):
        body = _get(server, ".gitignore")
        assert "<pre>" in body
        assert "__pycache__/" in body

    def test_known_filename_served_as_text(self, server):
        body = _get(server, "Makefile")
        assert "<pre>" in body

    def test_binary_with_text_ext_does_not_crash(self, server):
        """A .txt with invalid UTF-8 should fall back to binary serving."""
        with urlopen(f"{server}/broken.txt") as resp:
            assert resp.status == 200


# ---------------------------------------------------------------------------
# serve_directory_listing integration tests
# ---------------------------------------------------------------------------

class TestServeDirectoryListing:
    def test_root_lists_files(self, server):
        body = _get(server, "")
        assert "hello.py" in body
        assert "basic.md" in body

    def test_root_lists_subdirs(self, server):
        body = _get(server, "")
        assert "docs/" in body

    def test_subdir_with_readme_renders_markdown(self, server):
        body = _get(server, "docs")
        assert "<h1>Docs</h1>" in body

    def test_subdir_without_readme_shows_listing(self, server):
        body = _get(server, "empty_dir")
        assert "Directory" in body

    def test_subdir_listing_has_parent_link(self, server):
        body = _get(server, "empty_dir")
        assert ".." in body


# ---------------------------------------------------------------------------
# do_GET routing integration tests
# ---------------------------------------------------------------------------

class TestRouting:
    def test_markdown_routed_to_renderer(self, server):
        body = _get(server, "plain.md")
        assert "<h1>Just a heading</h1>" in body

    def test_binary_file_served_raw(self, server):
        with urlopen(f"{server}/image.png") as resp:
            data = resp.read()
            assert data[:4] == b"\x89PNG"

    def test_missing_file_returns_404(self, server):
        with pytest.raises(HTTPError) as exc_info:
            _get(server, "nonexistent.txt")
        assert exc_info.value.code == 404

    def test_url_encoded_path(self, server):
        body = _get(server, "hello.py")
        assert "<pre>" in body
