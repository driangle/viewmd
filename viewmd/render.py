"""HTML page rendering functions for viewmd."""
import html as html_mod


def _frontmatter_table(frontmatter):
    """Render frontmatter dict as an HTML table snippet."""
    if not frontmatter:
        return ""
    rows = ""
    for key, value in frontmatter.items():
        ek, ev = html_mod.escape(key), html_mod.escape(value)
        rows += (
            f'<tr><td class="fm-key">{ek}</td>'
            f'<td>{ev}</td></tr>\n')
    return f'<div class="frontmatter"><table>{rows}</table></div>\n'


def render_markdown_page(file_name, frontmatter, body_html, base_url):
    """Return a complete HTML page for rendered markdown content."""
    fm = _frontmatter_table(frontmatter)
    return f"""<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <base href="{base_url}">
    <title>{file_name}</title>
    <style>
        body {{
            max-width: 800px; margin: 40px auto; padding: 0 20px;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI',
                         Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6; color: #333;
        }}
        pre {{
            background: #f4f4f4; border: 1px solid #ddd;
            border-radius: 4px; padding: 12px; overflow-x: auto;
        }}
        code {{
            background: #f4f4f4; padding: 2px 6px;
            border-radius: 3px; font-family: 'Courier New', monospace;
        }}
        pre code {{ background: none; padding: 0; }}
        table {{
            border-collapse: collapse; width: 100%; margin: 20px 0;
        }}
        th, td {{
            border: 1px solid #ddd; padding: 8px 12px; text-align: left;
        }}
        th {{ background: #f4f4f4; }}
        a {{ color: #0066cc; text-decoration: none; }}
        a:hover {{ text-decoration: underline; }}
        blockquote {{
            border-left: 4px solid #ddd; margin: 0;
            padding-left: 20px; color: #666;
        }}
        img {{ max-width: 100%; height: auto; }}
        .frontmatter {{
            background: #f8f9fa; border: 1px solid #e1e4e8;
            border-radius: 6px; padding: 4px 12px;
            margin-bottom: 24px; font-size: 0.85em; color: #586069;
        }}
        .frontmatter table {{ width: auto; margin: 8px 0; border: none; }}
        .frontmatter td {{ border: none; padding: 2px 12px 2px 0; }}
        .frontmatter .fm-key {{ font-weight: 600; color: #444; }}
    </style>
</head>
<body>
    {fm}{body_html}
</body>
</html>"""


def render_text_page(file_name, escaped_content):
    """Return a complete HTML page for a text file."""
    return f"""<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{file_name}</title>
    <style>
        body {{
            max-width: 1000px; margin: 20px auto; padding: 0 20px;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI',
                         Roboto, sans-serif;
            background: #f6f8fa;
        }}
        .header {{
            background: white; border: 1px solid #d0d7de;
            border-radius: 6px 6px 0 0; padding: 12px 16px;
            font-weight: 600; border-bottom: 1px solid #d0d7de;
        }}
        .content {{
            background: white; border: 1px solid #d0d7de;
            border-top: none; border-radius: 0 0 6px 6px;
            padding: 16px; overflow-x: auto;
        }}
        pre {{
            margin: 0;
            font-family: 'SF Mono', 'Monaco', 'Inconsolata',
                         'Fira Mono', 'Courier New', monospace;
            font-size: 12px; line-height: 1.5;
            white-space: pre; word-wrap: normal;
        }}
    </style>
</head>
<body>
    <div class="header">{file_name}</div>
    <div class="content"><pre>{escaped_content}</pre></div>
</body>
</html>"""


def render_directory_page(display_path, parent_href, items):
    """Return a complete HTML page for a directory listing.

    Args:
        display_path: Path string for the heading.
        parent_href: Href for parent link, or None for root.
        items: List of (name, href, is_dir) tuples.
    """
    body_items = ""
    if parent_href is not None:
        body_items += (
            f'<li><a href="/{parent_href}" class="dir">..</a></li>\n')
    for name, href, is_dir in items:
        css = "dir" if is_dir else "file"
        label = f"{name}/" if is_dir else name
        body_items += (
            f'<li><a href="/{href}" class="{css}">{label}</a></li>\n')
    return f"""<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Directory listing</title>
    <style>
        body {{
            max-width: 800px; margin: 40px auto; padding: 0 20px;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI',
                         Roboto, sans-serif;
        }}
        ul {{ list-style: none; padding: 0; }}
        li {{ margin: 8px 0; }}
        a {{ text-decoration: none; color: #0066cc; }}
        a:hover {{ text-decoration: underline; }}
    </style>
</head>
<body>
    <h1>Directory: /{display_path}</h1>
    <ul>
{body_items}    </ul>
</body>
</html>"""
