package wrapstore

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWrapStore(t *testing.T) {
	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "wrapstore_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a valid wrap file
	wrap1 := `
apiVersion: krapper.kubotal.io/v1alpha1
kind: Wrap
name: test-wrap-1
version: v1
menuMode: grid
source:
  apiVersion: v1
  kind: Pod
`
	err = os.WriteFile(filepath.Join(tmpDir, "wrap1.yaml"), []byte(wrap1), 0644)
	if err != nil {
		t.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	ws, err := New(tmpDir, logger)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Verify initial load
	catalog := ws.GetCatalog()
	if len(catalog.Wraps) != 1 {
		t.Errorf("Expected 1 wrap, got %d", len(catalog.Wraps))
	}
	if catalog.Wraps[0].Name != "test-wrap-1" {
		t.Errorf("Expected wrap name test-wrap-1, got %s", catalog.Wraps[0].Name)
	}

	w := ws.GetWrap("test-wrap-1")
	if w.Name != "test-wrap-1" {
		t.Errorf("Expected wrap name test-wrap-1, got %s", w.Name)
	}

	// Test Watcher: Add a new file
	wrap2 := `
apiVersion: krapper.kubotal.io/v1alpha1
kind: Wrap
name: test-wrap-2
version: v1
menuMode: subMenu
source:
  apiVersion: v1
  kind: Service
`
	err = os.WriteFile(filepath.Join(tmpDir, "wrap2.yaml"), []byte(wrap2), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Give watcher some time
	time.Sleep(500 * time.Millisecond)

	catalog = ws.GetCatalog()
	if len(catalog.Wraps) != 2 {
		t.Errorf("Expected 2 wraps, got %d", len(catalog.Wraps))
	}

	w2 := ws.GetWrap("test-wrap-2")
	if w2.Name != "test-wrap-2" {
		t.Errorf("Expected wrap name test-wrap-2, got %s", w2.Name)
	}

	// Test Watcher: Update file
	wrap1Updated := `
apiVersion: krapper.kubotal.io/v1alpha1
kind: Wrap
name: test-wrap-1
version: v2
menuMode: grid
source:
  apiVersion: v1
  kind: Pod
`
	err = os.WriteFile(filepath.Join(tmpDir, "wrap1.yaml"), []byte(wrap1Updated), 0644)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)
	w = ws.GetWrap("test-wrap-1")
	if w.Version != "v2" {
		t.Errorf("Expected version v2, got %s", w.Version)
	}

	// Test Watcher: Remove file
	err = os.Remove(filepath.Join(tmpDir, "wrap2.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)
	catalog = ws.GetCatalog()
	if len(catalog.Wraps) != 1 {
		t.Errorf("Expected 1 wrap after deletion, got %d", len(catalog.Wraps))
	}

	// Test Watcher: New Directory
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	// sleep to allow watcher to add dir
	time.Sleep(500 * time.Millisecond)

	wrap3 := `
apiVersion: krapper.kubotal.io/v1alpha1
kind: Wrap
name: test-wrap-3
version: v1
menuMode: grid
source:
  apiVersion: v1
  kind: ConfigMap
`
	err = os.WriteFile(filepath.Join(subDir, "wrap3.yaml"), []byte(wrap3), 0644)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)
	catalog = ws.GetCatalog()
	if len(catalog.Wraps) != 2 {
		t.Errorf("Expected 2 wraps (1 + 1 in subdir), got %d", len(catalog.Wraps))
	}
}
