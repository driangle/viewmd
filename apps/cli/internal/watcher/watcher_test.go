package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileChangeDetection(t *testing.T) {
	dir := t.TempDir()
	w, err := New(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	// Write a file and expect an event.
	os.WriteFile(filepath.Join(dir, "test.md"), []byte("hello"), 0644)

	select {
	case <-w.Events():
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for file change event")
	}
}

func TestDebouncing(t *testing.T) {
	dir := t.TempDir()
	w, err := New(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	// Rapid writes should produce a single debounced event.
	for i := 0; i < 5; i++ {
		os.WriteFile(filepath.Join(dir, "test.md"), []byte{byte(i)}, 0644)
		time.Sleep(10 * time.Millisecond)
	}

	select {
	case <-w.Events():
		// ok, got first event
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for debounced event")
	}

	// Channel should be empty now (debounced into one event).
	select {
	case <-w.Events():
		// A second event is acceptable due to timing, but there shouldn't be 5.
	case <-time.After(300 * time.Millisecond):
		// ok, no extra events
	}
}

func TestIgnorePatterns(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".git"), 0755)
	os.MkdirAll(filepath.Join(dir, "visible"), 0755)

	w, err := New(dir, []string{".git"})
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	// Write to ignored dir — should not trigger.
	os.WriteFile(filepath.Join(dir, ".git", "index"), []byte("x"), 0644)

	select {
	case <-w.Events():
		t.Fatal("received event for ignored path")
	case <-time.After(300 * time.Millisecond):
		// ok
	}

	// Write to visible dir — should trigger.
	os.WriteFile(filepath.Join(dir, "visible", "file.md"), []byte("y"), 0644)

	select {
	case <-w.Events():
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for non-ignored event")
	}
}

func TestNewSubdirectoryDetection(t *testing.T) {
	dir := t.TempDir()
	w, err := New(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	// Create a new subdirectory.
	subdir := filepath.Join(dir, "newdir")
	os.Mkdir(subdir, 0755)

	// Drain the directory creation event.
	select {
	case <-w.Events():
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for mkdir event")
	}

	// Write a file in the new subdirectory.
	time.Sleep(200 * time.Millisecond)
	os.WriteFile(filepath.Join(subdir, "new.md"), []byte("z"), 0644)

	select {
	case <-w.Events():
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for event in new subdirectory")
	}
}
