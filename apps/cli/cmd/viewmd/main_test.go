package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"testing"
)

func TestPrintBanner(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printBanner(8000)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	checks := []string{
		"Markdown Server v0.1.0",
		"http://localhost:8000",
		"Markdown rendering",
		"Text file viewer",
		"Directory browsing",
		"Ctrl+C",
	}
	for _, s := range checks {
		if !strings.Contains(output, s) {
			t.Errorf("banner missing %q", s)
		}
	}
}

func TestPortInUseMessage(t *testing.T) {
	// Occupy a port to test the error path.
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port

	// Try binding to the same port — should fail.
	_, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err == nil {
		t.Fatal("expected port-in-use error")
	}
}
