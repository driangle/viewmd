#!/usr/bin/env python3
"""
viewmd - Simple HTTP server for viewing Markdown files in your browser.
Usage: viewmd [port]
"""
import sys
import html
from http.server import HTTPServer, SimpleHTTPRequestHandler
from pathlib import Path
import urllib.parse
from datetime import datetime

VERSION = "0.1.5"

try:
    import markdown
except ImportError:
    print("Error: 'markdown' package not found.")
    print("Install it with: pip3 install markdown")
    sys.exit(1)

from .render import (  # pylint: disable=wrong-import-position
    render_markdown_page, render_text_page, render_directory_page)


def parse_frontmatter(content):
    """Extract frontmatter dict and remaining content from markdown text."""
    if not content.startswith('---'):
        return None, content
    parts = content.split('---', 2)
    if len(parts) < 3:
        return None, content
    frontmatter = {}
    for line in parts[1].strip().splitlines():
        if ':' in line:
            key, _, value = line.partition(':')
            frontmatter[key.strip()] = value.strip()
    return frontmatter, parts[2]


class MarkdownHandler(SimpleHTTPRequestHandler):
    """HTTP request handler that renders Markdown and text files."""

    TEXT_EXTENSIONS = {
        '.txt', '.log', '.json', '.xml', '.yaml', '.yml',
        '.toml', '.ini', '.cfg', '.conf', '.sh', '.bash',
        '.zsh', '.fish', '.py', '.js', '.ts', '.jsx', '.tsx',
        '.java', '.c', '.cpp', '.h', '.hpp', '.cs', '.go',
        '.rs', '.rb', '.php', '.swift', '.kt', '.sql', '.html',
        '.css', '.scss', '.sass', '.less', '.vue', '.svelte',
        '.r', '.m', '.scala', '.pl', '.lua', '.vim', '.el',
        '.clj', '.ex', '.exs', '.dockerfile', '.env',
        '.gitignore', '.gitattributes', '.editorconfig',
        '.eslintrc', '.prettierrc', '.babelrc',
    }

    TEXT_FILENAMES = {
        'makefile', 'dockerfile', 'gemfile', 'rakefile',
        'procfile', 'jenkinsfile', 'license', 'readme',
        'changelog', 'authors', 'contributors', 'codeowners',
    }

    def do_GET(self):
        """Route requests to the appropriate handler."""
        parsed = urllib.parse.urlparse(self.path)
        path = urllib.parse.unquote(parsed.path.lstrip('/'))
        timestamp = datetime.now().strftime('%H:%M:%S')
        print(f"[{timestamp}] Request: {path or '/'}")

        if not path:
            self.serve_directory_listing()
            return

        file_path = Path(path)

        if file_path.is_dir():
            readme = file_path / "README.md"
            if readme.is_file():
                self.serve_markdown(readme)
            else:
                self.serve_directory_listing()
            return

        if not file_path.is_file():
            super().do_GET()
            return

        if file_path.suffix.lower() in ('.md', '.markdown'):
            self.serve_markdown(file_path)
        elif self.is_text_file(file_path):
            self.serve_text_file(file_path)
        else:
            super().do_GET()

    def is_text_file(self, file_path):
        """Check if a file should be displayed as text."""
        if file_path.suffix.lower() in self.TEXT_EXTENSIONS:
            return True
        if file_path.name.lower() in self.TEXT_FILENAMES:
            return True
        return file_path.name.startswith('.')

    def serve_markdown(self, file_path):
        """Serve a markdown file as rendered HTML."""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            frontmatter, body = parse_frontmatter(content)
            md = markdown.Markdown(
                extensions=['fenced_code', 'tables', 'nl2br'])
            body_html = md.convert(body)
            base = file_path.parent
            base_url = f"/{base}/" if str(base) != '.' else "/"
            page = render_markdown_page(
                file_path.name, frontmatter, body_html, base_url)
            self._send_html(page)
        except Exception as e:  # pylint: disable=broad-exception-caught
            self.send_error(500, f"Error rendering markdown: {e}")

    def serve_text_file(self, file_path):
        """Serve a text file with HTML formatting."""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            page = render_text_page(file_path.name, html.escape(content))
            self._send_html(page)
        except UnicodeDecodeError:
            super().do_GET()
        except Exception as e:  # pylint: disable=broad-exception-caught
            self.send_error(500, f"Error displaying file: {e}")

    def serve_directory_listing(self):
        """Serve an HTML directory listing."""
        try:
            parsed = urllib.parse.urlparse(self.path)
            path = urllib.parse.unquote(parsed.path.lstrip('/')) or '.'
            dir_path = Path(path)
            if not dir_path.is_dir():
                self.send_error(404, "Directory not found")
                return
            entries = sorted(
                dir_path.iterdir(),
                key=lambda x: (not x.is_dir(), x.name.lower()))
            items = [
                (e.name, str(e.relative_to('.')), e.is_dir())
                for e in entries]
            parent = str(dir_path.parent) if path != '.' else None
            page = render_directory_page(path, parent, items)
            self._send_html(page)
        except Exception as e:  # pylint: disable=broad-exception-caught
            self.send_error(500, f"Error listing directory: {e}")

    def _send_html(self, content):
        """Send an HTML response."""
        self.send_response(200)
        self.send_header('Content-type', 'text/html; charset=utf-8')
        self.end_headers()
        self.wfile.write(content.encode('utf-8'))


def main():
    """Start the viewmd HTTP server."""
    port = int(sys.argv[1]) if len(sys.argv) > 1 else 8000
    try:
        server = HTTPServer(('', port), MarkdownHandler)
    except OSError:
        print(f"Error: Port {port} is already in use.")
        sys.exit(1)
    print("=" * 60)
    print(f"Markdown Server v{VERSION}")
    print("=" * 60)
    print(f"Server: http://localhost:{port}")
    print("Features:")
    print("  - Markdown rendering (.md, .markdown)")
    print("  - Text file viewer (.py, .js, .gitignore, etc.)")
    print("  - Directory browsing")
    print("=" * 60)
    print("Press Ctrl+C to stop\n")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\nShutting down...")
        server.shutdown()


if __name__ == '__main__':
    main()
