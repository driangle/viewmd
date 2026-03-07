// Package frontmatter parses YAML-like frontmatter from markdown content.
package frontmatter

import "strings"

// Parse extracts frontmatter key-value pairs from content delimited by "---".
// Returns (nil, content) when no valid frontmatter is found.
// Returns (map, body) when frontmatter is present, where body is everything
// after the closing "---" delimiter.
func Parse(content string) (map[string]string, string) {
	if !strings.HasPrefix(content, "---") {
		return nil, content
	}

	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return nil, content
	}

	meta := make(map[string]string)
	block := strings.TrimSpace(parts[1])
	if block != "" {
		for _, line := range strings.Split(block, "\n") {
			if idx := strings.Index(line, ":"); idx >= 0 {
				key := strings.TrimSpace(line[:idx])
				value := strings.TrimSpace(line[idx+1:])
				meta[key] = value
			}
		}
	}

	return meta, parts[2]
}
