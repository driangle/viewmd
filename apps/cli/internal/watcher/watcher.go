// Package watcher provides recursive filesystem watching with debouncing.
package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

const debounceInterval = 100 * time.Millisecond

// Watcher monitors a directory tree for file changes and sends
// debounced notifications on its Events channel.
type Watcher struct {
	fsw    *fsnotify.Watcher
	events chan struct{}
	done   chan struct{}
	ignore []string
}

// New creates a Watcher that recursively monitors root.
// ignore is a list of glob patterns for paths to skip.
func New(root string, ignore []string) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		fsw:    fsw,
		events: make(chan struct{}, 1),
		done:   make(chan struct{}),
		ignore: ignore,
	}

	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		name := d.Name()
		rel, _ := filepath.Rel(root, path)
		if rel == "." {
			return fsw.Add(path)
		}
		if w.isIgnored(name, rel) {
			return filepath.SkipDir
		}
		return fsw.Add(path)
	})
	if err != nil {
		fsw.Close()
		return nil, err
	}

	go w.run(root)
	return w, nil
}

// Events returns a channel that receives a signal on each debounced change.
func (w *Watcher) Events() <-chan struct{} {
	return w.events
}

// Close stops the watcher.
func (w *Watcher) Close() {
	w.fsw.Close()
	<-w.done
}

func (w *Watcher) run(root string) {
	defer close(w.done)

	var timer *time.Timer
	var timerC <-chan time.Time

	for {
		select {
		case ev, ok := <-w.fsw.Events:
			if !ok {
				return
			}
			rel, err := filepath.Rel(root, ev.Name)
			if err != nil {
				continue
			}
			name := filepath.Base(rel)
			if w.isIgnored(name, rel) {
				continue
			}
			// Auto-add new directories.
			if ev.Has(fsnotify.Create) {
				if info, err := os.Stat(ev.Name); err == nil && info.IsDir() {
					w.fsw.Add(ev.Name)
				}
			}
			if timer == nil {
				timer = time.NewTimer(debounceInterval)
				timerC = timer.C
			} else {
				timer.Reset(debounceInterval)
			}

		case <-timerC:
			timer = nil
			timerC = nil
			// Non-blocking send.
			select {
			case w.events <- struct{}{}:
			default:
			}

		case _, ok := <-w.fsw.Errors:
			if !ok {
				return
			}
		}
	}
}

// isIgnored returns true if the given name or relative path matches
// any of the ignore patterns.
func (w *Watcher) isIgnored(name, relPath string) bool {
	for _, p := range w.ignore {
		if strings.Contains(p, "/") {
			if matched, _ := filepath.Match(p, relPath); matched {
				return true
			}
		} else {
			if matched, _ := filepath.Match(p, name); matched {
				return true
			}
		}
	}
	return false
}
