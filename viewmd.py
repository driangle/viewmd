#!/usr/bin/env python3
"""
viewmd - Simple HTTP server for viewing Markdown files in your browser.
Usage: viewmd [port]
"""
import os
import sys
import html
from http.server import HTTPServer, SimpleHTTPRequestHandler
from pathlib import Path
import urllib.parse
from datetime import datetime

VERSION = "0.1.0"

try:
    import markdown
except ImportError:
    print("Error: 'markdown' package not found.")
    print("Install it with: pip3 install markdown")
    sys.exit(1)


class MarkdownHandler(SimpleHTTPRequestHandler):
    # Text file extensions to display in browser
    TEXT_EXTENSIONS = {
        '.txt', '.log', '.json', '.xml', '.yaml', '.yml', '.toml', '.ini', '.cfg', '.conf',
        '.sh', '.bash', '.zsh', '.fish', '.py', '.js', '.ts', '.jsx', '.tsx', '.java', '.c',
        '.cpp', '.h', '.hpp', '.cs', '.go', '.rs', '.rb', '.php', '.swift', '.kt', '.sql',
        '.html', '.css', '.scss', '.sass', '.less', '.vue', '.svelte', '.r', '.m', '.scala',
        '.pl', '.lua', '.vim', '.el', '.clj', '.ex', '.exs', '.dockerfile', '.env', '.gitignore',
        '.gitattributes', '.editorconfig', '.eslintrc', '.prettierrc', '.babelrc'
    }

    # Files without extensions that are typically text
    TEXT_FILENAMES = {
        'makefile', 'dockerfile', 'gemfile', 'rakefile', 'procfile', 'jenkinsfile',
        'license', 'readme', 'changelog', 'authors', 'contributors', 'codeowners'
    }

    def do_GET(self):
        # Parse the URL and remove query parameters
        parsed_path = urllib.parse.urlparse(self.path)
        path = parsed_path.path.lstrip('/')

        # Decode URL encoding
        path = urllib.parse.unquote(path)

        print(f"[{datetime.now().strftime('%H:%M:%S')}] Request: {path or '/'}")

        # Handle root directory
        if not path:
            print(f"  → Serving directory listing")
            self.serve_directory_listing()
            return

        # Get the file path
        file_path = Path(path)

        # Check if it's a directory
        if file_path.is_dir():
            # Check for README.md in the directory
            readme_path = file_path / "README.md"
            if readme_path.is_file():
                print(f"  → Rendering markdown: {readme_path}")
                self.serve_markdown(readme_path)
            else:
                print(f"  → Serving directory listing")
                self.serve_directory_listing()
            return

        # Check if file exists
        if not file_path.is_file():
            print(f"  → File not found, using default handler")
            super().do_GET()
            return

        # Check if it's a markdown file
        if file_path.suffix.lower() in ['.md', '.markdown']:
            print(f"  → Rendering markdown: {file_path}")
            self.serve_markdown(file_path)
        # Check if it's a text file
        elif self.is_text_file(file_path):
            print(f"  → Displaying text file: {file_path}")
            self.serve_text_file(file_path)
        else:
            # Serve other files normally (images, PDFs, etc.)
            print(f"  → Serving binary file: {file_path}")
            super().do_GET()

    def is_text_file(self, file_path):
        """Check if a file should be displayed as text."""
        filename = file_path.name

        # Check extension
        if file_path.suffix.lower() in self.TEXT_EXTENSIONS:
            print(f"  → Detected as text (extension: {file_path.suffix})")
            return True

        # Check filename (for files without extensions)
        if filename.lower() in self.TEXT_FILENAMES:
            print(f"  → Detected as text (known filename: {filename})")
            return True

        # Check if filename starts with a dot (hidden config files like .gitignore, .env, etc.)
        # In pathlib, .gitignore has suffix='.gitignore' and stem='', so we check if name starts with '.'
        if filename.startswith('.'):
            print(f"  → Detected as text (dotfile: {filename})")
            return True

        print(f"  → Not detected as text file")
        return False

    def serve_markdown(self, file_path):
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()

            # Convert markdown to HTML
            md = markdown.Markdown(extensions=['fenced_code', 'tables', 'nl2br'])
            html_content = md.convert(content)

            # Calculate the base URL for relative links
            # This ensures links in markdown are relative to the file's directory
            base_path = file_path.parent if file_path.parent != Path('.') else Path('/')
            base_url = f"/{base_path}/" if str(base_path) != '.' else "/"

            # Wrap in a nice HTML template
            html = f"""<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <base href="{base_url}">
    <title>{file_path.name}</title>
    <style>
        body {{
            max-width: 800px;
            margin: 40px auto;
            padding: 0 20px;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
        }}
        pre {{
            background: #f4f4f4;
            border: 1px solid #ddd;
            border-radius: 4px;
            padding: 12px;
            overflow-x: auto;
        }}
        code {{
            background: #f4f4f4;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
        }}
        pre code {{
            background: none;
            padding: 0;
        }}
        table {{
            border-collapse: collapse;
            width: 100%;
            margin: 20px 0;
        }}
        th, td {{
            border: 1px solid #ddd;
            padding: 8px 12px;
            text-align: left;
        }}
        th {{
            background: #f4f4f4;
        }}
        a {{
            color: #0066cc;
            text-decoration: none;
        }}
        a:hover {{
            text-decoration: underline;
        }}
        blockquote {{
            border-left: 4px solid #ddd;
            margin: 0;
            padding-left: 20px;
            color: #666;
        }}
        img {{
            max-width: 100%;
            height: auto;
        }}
    </style>
</head>
<body>
    {html_content}
</body>
</html>"""

            # Send response
            self.send_response(200)
            self.send_header('Content-type', 'text/html; charset=utf-8')
            self.end_headers()
            self.wfile.write(html.encode('utf-8'))

        except Exception as e:
            self.send_error(500, f"Error rendering markdown: {str(e)}")

    def serve_text_file(self, file_path):
        """Serve a text file with HTML formatting."""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()

            # Escape HTML characters
            escaped_content = html.escape(content)

            # Wrap in HTML template
            html_output = f"""<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{file_path.name}</title>
    <style>
        body {{
            max-width: 1000px;
            margin: 20px auto;
            padding: 0 20px;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #f6f8fa;
        }}
        .header {{
            background: white;
            border: 1px solid #d0d7de;
            border-radius: 6px 6px 0 0;
            padding: 12px 16px;
            font-weight: 600;
            border-bottom: 1px solid #d0d7de;
        }}
        .content {{
            background: white;
            border: 1px solid #d0d7de;
            border-top: none;
            border-radius: 0 0 6px 6px;
            padding: 16px;
            overflow-x: auto;
        }}
        pre {{
            margin: 0;
            font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Fira Mono', 'Courier New', monospace;
            font-size: 12px;
            line-height: 1.5;
            white-space: pre;
            word-wrap: normal;
        }}
    </style>
</head>
<body>
    <div class="header">{file_path.name}</div>
    <div class="content">
        <pre>{escaped_content}</pre>
    </div>
</body>
</html>"""

            # Send response
            self.send_response(200)
            self.send_header('Content-type', 'text/html; charset=utf-8')
            self.end_headers()
            self.wfile.write(html_output.encode('utf-8'))

        except UnicodeDecodeError:
            # If it's not valid UTF-8, serve it normally (might be binary)
            super().do_GET()
        except Exception as e:
            self.send_error(500, f"Error displaying file: {str(e)}")

    def serve_directory_listing(self):
        try:
            parsed_path = urllib.parse.urlparse(self.path)
            path = parsed_path.path.lstrip('/') or '.'
            path = urllib.parse.unquote(path)

            dir_path = Path(path)
            if not dir_path.is_dir():
                self.send_error(404, "Directory not found")
                return

            items = sorted(dir_path.iterdir(), key=lambda x: (not x.is_dir(), x.name.lower()))

            html = """<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Directory listing</title>
    <style>
        body {
            max-width: 800px;
            margin: 40px auto;
            padding: 0 20px;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        }
        ul { list-style: none; padding: 0; }
        li { margin: 8px 0; }
        a { text-decoration: none; color: #0066cc; }
        a:hover { text-decoration: underline; }
        .dir::before { content: '📁 '; }
        .file::before { content: '📄 '; }
    </style>
</head>
<body>
"""
            html += f"<h1>Directory: /{path}</h1>\n<ul>\n"

            if path != '.':
                parent = dir_path.parent
                html += f'<li><a href="/{parent}" class="dir">..</a></li>\n'

            for item in items:
                rel_path = item.relative_to('.')
                if item.is_dir():
                    html += f'<li><a href="/{rel_path}" class="dir">{item.name}/</a></li>\n'
                else:
                    html += f'<li><a href="/{rel_path}" class="file">{item.name}</a></li>\n'

            html += "</ul>\n</body>\n</html>"

            self.send_response(200)
            self.send_header('Content-type', 'text/html; charset=utf-8')
            self.end_headers()
            self.wfile.write(html.encode('utf-8'))

        except Exception as e:
            self.send_error(500, f"Error listing directory: {str(e)}")


def main():
    port = int(sys.argv[1]) if len(sys.argv) > 1 else 8000

    server = HTTPServer(('', port), MarkdownHandler)
    print(f"=" * 60)
    print(f"Markdown Server v{VERSION}")
    print(f"=" * 60)
    print(f"Server: http://localhost:{port}")
    print(f"Features:")
    print(f"  - Markdown rendering (.md, .markdown)")
    print(f"  - Text file viewer (.py, .js, .gitignore, etc.)")
    print(f"  - Directory browsing")
    print(f"=" * 60)
    print("Press Ctrl+C to stop\n")

    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\nShutting down...")
        server.shutdown()


if __name__ == '__main__':
    main()
