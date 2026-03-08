package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestServeSSE_SendsReloadEvent(t *testing.T) {
	ch := make(chan struct{}, 1)
	h := New(".")
	h.WatchMode = true
	h.SetWatchEvents(ch)

	srv := httptest.NewServer(h)
	defer srv.Close()

	// Send an event before connecting so the client sees it immediately.
	ch <- struct{}{}

	resp, err := http.Get(srv.URL + "/-/events")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "text/event-stream") {
		t.Fatalf("expected text/event-stream, got %q", ct)
	}

	buf := make([]byte, 256)
	n, err := resp.Body.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	body := string(buf[:n])
	if !strings.Contains(body, "data: reload") {
		t.Fatalf("expected reload event, got %q", body)
	}
}

func TestServeSSE_DisabledWhenNotWatchMode(t *testing.T) {
	h := New(".")
	h.WatchMode = false

	srv := httptest.NewServer(h)
	defer srv.Close()

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(srv.URL + "/-/events")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Should get 404 when watch mode is off.
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}
