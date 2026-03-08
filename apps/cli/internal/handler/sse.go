package handler

import (
	"fmt"
	"net/http"
)

// SetWatchEvents configures the handler for watch mode by providing
// the channel that signals file changes.
func (h *Handler) SetWatchEvents(ch <-chan struct{}) {
	h.watchEvents = ch
}

// serveSSE streams Server-Sent Events to the client, sending a
// "reload" message each time a file change is detected.
func (h *Handler) serveSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case _, ok := <-h.watchEvents:
			if !ok {
				return
			}
			fmt.Fprintf(w, "data: reload\n\n")
			flusher.Flush()
		}
	}
}
