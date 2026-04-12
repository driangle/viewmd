// Package wikilink converts wiki-style [[links]] in markdown text
// to standard markdown links.
package wikilink

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

// wikiLinkRe matches [[target]] or [[target|display text]].
var wikiLinkRe = regexp.MustCompile(`\[\[([^\]\|]+?)(?:\|([^\]]+?))?\]\]`)

// Resolve replaces wiki-style [[links]] in markdown text with standard
// markdown links. Links are resolved relative to baseURL.
//
// Supported forms:
//
//	[[page-name]]           → [page-name](page-name.md)
//	[[page-name|Display]]   → [Display](page-name.md)
//	[[sub/page]]            → [sub/page](sub/page.md)
//	[[file.md]]             → [file.md](file.md)     (no double extension)
func Resolve(text string, baseURL string) string {
	return wikiLinkRe.ReplaceAllStringFunc(text, func(match string) string {
		parts := wikiLinkRe.FindStringSubmatch(match)
		target := strings.TrimSpace(parts[1])
		display := strings.TrimSpace(parts[2])

		if target == "" {
			return match
		}

		if display == "" {
			display = target
		}

		href := resolveHref(target, baseURL)
		return fmt.Sprintf("[%s](%s)", display, href)
	})
}

// resolveHref builds the link href from a wiki target and base URL.
// Appends .md if the target has no markdown extension.
func resolveHref(target string, baseURL string) string {
	ext := strings.ToLower(path.Ext(target))
	if ext != ".md" && ext != ".markdown" {
		target = target + ".md"
	}
	// Ensure baseURL ends with /
	if baseURL != "" && !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return baseURL + target
}
