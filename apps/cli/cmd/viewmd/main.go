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
	"github.com/driangle/viewmd/apps/cli/internal/watcher"
)

var version = "0.1.0"

func main() {
	render.Version = version

	autoReadme := flag.Bool("auto-readme", false, "Auto-render README.md in directories")
	showAll := flag.Bool("show-all", false, "Show all files, not just Markdown (shorthand: -a)")
	flag.BoolVar(showAll, "a", false, "Show all files, not just Markdown")
	watch := flag.Bool("watch", false, "Watch for file changes and auto-reload browser (shorthand: -w)")
	flag.BoolVar(watch, "w", false, "Watch for file changes and auto-reload browser")
	ignoreFlag := flag.String("ignore", "", "Comma-separated ignore patterns (glob syntax)")
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

	ignorePatterns := buildIgnorePatterns(*ignoreFlag, root)

	printBanner(port, root)
	printArgs(port, *autoReadme, *showAll, *watch, ignorePatterns)

	h := handler.New(root)
	h.AutoReadme = *autoReadme
	h.ShowAll = *showAll || loadShowAllFromConfig(root)
	h.IgnorePatterns = ignorePatterns

	if *watch {
		w, err := watcher.New(root, ignorePatterns)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to start watcher: %v\n", err)
			os.Exit(1)
		}
		defer w.Close()
		h.WatchMode = true
		h.SetWatchEvents(w.Events())
		render.WatchMode = true
	}

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

func printArgs(port int, autoReadme, showAll, watch bool, ignorePatterns []string) {
	fmt.Fprintf(os.Stderr, "  port:        %d\n", port)
	fmt.Fprintf(os.Stderr, "  auto-readme: %v\n", autoReadme)
	fmt.Fprintf(os.Stderr, "  show-all:    %v\n", showAll)
	fmt.Fprintf(os.Stderr, "  watch:       %v\n", watch)
	if len(ignorePatterns) > 0 {
		fmt.Fprintf(os.Stderr, "  ignore:      %s\n", strings.Join(ignorePatterns, ", "))
	}
	fmt.Fprintln(os.Stderr)
}

// buildIgnorePatterns merges CLI and YAML ignore patterns.
// If no patterns from any source, defaults to [".git"].
func buildIgnorePatterns(flagVal, root string) []string {
	var patterns []string

	if flagVal != "" {
		for _, p := range strings.Split(flagVal, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				patterns = append(patterns, p)
			}
		}
	}

	patterns = append(patterns, loadIgnoreFromConfig(root)...)

	if len(patterns) == 0 {
		return []string{".git"}
	}
	return patterns
}

// loadIgnoreFromConfig reads the ignore list from .viewmd.yaml.
func loadIgnoreFromConfig(root string) []string {
	f, err := os.Open(filepath.Join(root, ".viewmd.yaml"))
	if err != nil {
		return nil
	}
	defer f.Close()

	var patterns []string
	inIgnore := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "ignore:" {
			inIgnore = true
			continue
		}

		if inIgnore {
			if strings.HasPrefix(trimmed, "- ") {
				val := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
				val = strings.Trim(val, `"'`)
				if val != "" {
					patterns = append(patterns, val)
				}
			} else {
				inIgnore = false
			}
		}
	}
	return patterns
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
