package wrapstore

import (
	"fmt"
	"io/fs"
	"krapper/internal/wrap"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/fsnotify.v1"
)

type CatalogItem struct {
	Name     string        `yaml:"name" json:"name"`
	Label    string        `yaml:"label" json:"label"`
	MenuMode wrap.MenuMode `yaml:"menuMode" json:"menuMode"`
}

type Catalog struct {
	Wraps []CatalogItem `yaml:"wraps" json:"wraps"`
}

type WrapStore interface {
	GetCatalog() *Catalog
	GetWrap(name string) *wrap.Wrap
}

type store struct {
	mu      sync.RWMutex
	wraps   map[string]*wrap.Wrap
	catalog *Catalog
	files   map[string]string // filePath -> wrapName
	watcher *fsnotify.Watcher
	logger  *slog.Logger
	baseDir string
}

func New(path string, logger *slog.Logger) (WrapStore, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	s := &store{
		wraps:   make(map[string]*wrap.Wrap),
		files:   make(map[string]string),
		logger:  logger,
		baseDir: absPath,
	}

	// 1. Initial Load
	if err := s.loadAll(); err != nil {
		return nil, err
	}
	s.rebuildCatalog()

	// 2. Setup Watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}
	s.watcher = watcher

	// Add all subdirectories to watcher
	if err := s.watchRecursive(absPath); err != nil {
		_ = watcher.Close()
		return nil, err
	}

	// 3. Start Watch Loop
	go s.watchLoop()

	return s, nil
}

func (s *store) GetCatalog() *Catalog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.catalog
}

func (s *store) GetWrap(name string) *wrap.Wrap {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.wraps[name]
}

func (s *store) loadAll() error {
	return filepath.WalkDir(s.baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), ".yaml") {
			s.updateFile(path)
		}
		return nil
	})
}

func (s *store) updateFile(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Load the file
	w, err := wrap.Load(path)
	if err != nil {
		s.logger.Warn("Failed to load wrap", "path", path, "error", err)
		return
	}
	if w == nil {
		s.logger.Warn("File is not a wrap, skipping", "path", path)
		return
	}

	// Check if this file was previously associated with a different wrap name (unlikely but possible if name changed in file)
	if oldName, ok := s.files[path]; ok {
		if oldName != w.Name {
			delete(s.wraps, oldName)
		}
	}

	s.wraps[w.Name] = w
	s.files[path] = w.Name
	s.logger.Info("Loaded wrap", "name", w.Name, "path", path)
	s.rebuildCatalog()
}

func (s *store) removeFile(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if name, ok := s.files[path]; ok {
		delete(s.wraps, name)
		delete(s.files, path)
		s.rebuildCatalog()
	}
}

func (s *store) rebuildCatalog() {
	catalog := &Catalog{
		Wraps: make([]CatalogItem, 0, len(s.wraps)),
	}

	for _, w := range s.wraps {
		catalog.Wraps = append(catalog.Wraps, CatalogItem{
			Name:     w.Name,
			Label:    w.Label,
			MenuMode: w.MenuMode,
		})
	}
	s.catalog = catalog
}

func (s *store) watchRecursive(path string) error {
	return filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return s.watcher.Add(p)
		}
		return nil
	})
}

func (s *store) watchLoop() {
	defer func() { _ = s.watcher.Close() }()

	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}

			// Handle new directories (Watcher doesn't recursively watch new dirs automatically)
			// But note: fsnotify events order for mkdir might vary.
			// Ideally we check if it is a directory on CREATE.
			if event.Op&fsnotify.Create == fsnotify.Create {
				stat, err := os.Stat(event.Name)
				if err == nil && stat.IsDir() {
					_ = s.watcher.Add(event.Name)
				}
			}

			if strings.HasSuffix(event.Name, ".yaml") {
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					s.updateFile(event.Name)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
					s.removeFile(event.Name)
				}
			}

		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
			s.logger.Error("watcher error", "error", err)
		}
	}
}
