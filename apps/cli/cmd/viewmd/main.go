package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/driangle/viewmd/apps/cli/internal/handler"
	"github.com/driangle/viewmd/apps/cli/internal/logging"
	"github.com/driangle/viewmd/apps/cli/internal/render"
)

var version = "0.2.0"

func main() {
	render.Version = version

	autoReadme := flag.Bool("auto-readme", false, "Auto-render README.md in directories")
	showAll := flag.Bool("show-all", false, "Show all files, not just Markdown (shorthand: -a)")
	flag.BoolVar(showAll, "a", false, "Show all files, not just Markdown")
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

	root, _ := os.Getwd()
	printBanner(port, root)
	printArgs(port, *autoReadme, *showAll)

	h := handler.New(root)
	h.AutoReadme = *autoReadme
	h.ShowAll = *showAll || loadShowAllFromConfig(root)
	srv := &http.Server{Handler: logging.Middleware(h)}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}()

	<-done
	fmt.Fprintln(os.Stderr, "\nShutting down.")
	srv.Shutdown(context.Background())
}

func printBanner(port int, root string) {
	fmt.Fprintf(os.Stderr, "viewmd v%s\n", version)
	fmt.Fprintf(os.Stderr, "Serving %s on http://localhost:%d\n", root, port)
}

func printArgs(port int, autoReadme, showAll bool) {
	fmt.Fprintf(os.Stderr, "  port:        %d\n", port)
	fmt.Fprintf(os.Stderr, "  auto-readme: %v\n", autoReadme)
	fmt.Fprintf(os.Stderr, "  show-all:    %v\n\n", showAll)
}

// loadShowAllFromConfig reads .viewmd.yaml from root and returns
// the show_all_files value. Returns false if the file doesn't exist
// or can't be parsed.
func loadShowAllFromConfig(root string) bool {
	f, err := os.Open(filepath.Join(root, ".viewmd.yaml"))
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "show_all_files:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "show_all_files:"))
			return val == "true"
		}
	}
	return false
}
