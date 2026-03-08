// Package frontmatter parses YAML-like frontmatter from markdown content.
package frontmatter

import "strings"

// KeyValue holds a single frontmatter key-value pair.
type KeyValue struct {
	Key   string
	Value string
}

// Parse extracts frontmatter key-value pairs from content delimited by "---".
// Returns (nil, content) when no valid frontmatter is found.
// Returns (pairs, body) when frontmatter is present, where body is everything
// after the closing "---" delimiter. Pairs preserve the original order.
func Parse(content string) ([]KeyValue, string) {
	if !strings.HasPrefix(content, "---") {
		return nil, content
	}

	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return nil, content
	}

	var pairs []KeyValue
	block := strings.TrimSpace(parts[1])
	if block != "" {
		for _, line := range strings.Split(block, "\n") {
			if idx := strings.Index(line, ":"); idx >= 0 {
				key := strings.TrimSpace(line[:idx])
				value := strings.TrimSpace(line[idx+1:])
				pairs = append(pairs, KeyValue{Key: key, Value: value})
			}
		}
	}

	if pairs == nil {
		pairs = []KeyValue{}
	}

	return pairs, parts[2]
}
