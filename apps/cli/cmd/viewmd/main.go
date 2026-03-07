package main

import (
	"context"
	"errors"
	"flag"
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

	autoReadme := flag.Bool("auto-readme", false, "Auto-render README.md in directories")
	ver := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *ver {
		fmt.Println("viewmd", version)
		return
	}

	port := 8000
	if flag.NArg() > 0 {
		p, err := strconv.Atoi(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid port %q\n", flag.Arg(0))
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
	h := handler.New(root)
	h.AutoReadme = *autoReadme
	srv := &http.Server{Handler: h}

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
