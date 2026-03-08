package logging

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// responseRecorder wraps http.ResponseWriter to capture the status code.
type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Middleware returns an HTTP handler that logs each request as:
//
//	HH:MM:SS METHOD /path STATUS duration
func Middleware(next http.Handler) http.Handler {
	color := ColorEnabled()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		elapsed := time.Since(start)
		ts := start.Format("15:04:05")

		fmt.Fprintf(os.Stderr, "%s %s %s %s %s\n",
			dim(ts, color),
			r.Method,
			r.URL.Path,
			statusColor(rec.status, color),
			dim(formatDuration(elapsed), color),
		)
	})
}

// formatDuration returns a human-friendly duration string.
func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%d\u00b5s", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}
