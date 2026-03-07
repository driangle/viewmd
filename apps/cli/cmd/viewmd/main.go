package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/driangle/viewmd/apps/cli/internal/handler"
	"github.com/driangle/viewmd/apps/cli/internal/render"
)

const version = "0.1.0"

func main() {
	render.Version = version

	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("viewmd", version)
		return
	}

	port := 8000
	if len(os.Args) > 1 {
		p, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid port %q\n", os.Args[1])
			os.Exit(1)
		}
		port = p
	}

	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Port %d is already in use.\n", port)
		os.Exit(1)
	}

	printBanner(port)

	root, _ := os.Getwd()
	srv := &http.Server{Handler: handler.New(root)}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}()

	<-done
	fmt.Println("\nShutting down...")
	srv.Shutdown(context.Background())
}

func printBanner(port int) {
	sep := "============================================================"
	fmt.Println(sep)
	fmt.Printf("Markdown Server v%s\n", version)
	fmt.Println(sep)
	fmt.Printf("Server: http://localhost:%d\n", port)
	fmt.Println("Features:")
	fmt.Println("  - Markdown rendering (.md, .markdown)")
	fmt.Println("  - Text file viewer (.py, .js, .gitignore, etc.)")
	fmt.Println("  - Directory browsing")
	fmt.Println(sep)
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()
}
