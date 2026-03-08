// Package logging provides request logging middleware for the viewmd server.
package logging

import (
	"fmt"
	"os"
)

// ANSI color codes for terminal output.
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorDim    = "\033[2m"
)

// ColorEnabled reports whether the terminal likely supports ANSI colors.
// Respects the NO_COLOR convention (https://no-color.org/) and checks
// that TERM is set (indicating a real terminal, not a pipe).
func ColorEnabled() bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}
	term := os.Getenv("TERM")
	return term != "" && term != "dumb"
}

// statusColor returns the ANSI-colored status code string.
func statusColor(code int, color bool) string {
	s := fmt.Sprintf("%d", code)
	if !color {
		return s
	}
	var c string
	switch {
	case code >= 500:
		c = colorRed
	case code >= 400:
		c = colorYellow
	default:
		c = colorGreen
	}
	return c + s + colorReset
}

// dim wraps text in dim ANSI codes when color is enabled.
func dim(s string, color bool) string {
	if !color {
		return s
	}
	return colorDim + s + colorReset
}
