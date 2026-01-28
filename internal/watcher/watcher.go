package watcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"dugong/internal/config"
	"dugong/internal/generator"
	"dugong/internal/parser"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	config  *config.Config
	generator *generator.Generator
	watcher *fsnotify.Watcher
	timers  map[string]*time.Timer
	timerMu sync.Mutex
}

func NewWatcher(cfg *config.Config, generator *generator.Generator) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &Watcher{
		config:  cfg,
		generator: generator,
		watcher: watcher,
		timers:  make(map[string]*time.Timer),
	}, nil
}

func (w *Watcher) Watch() error {
	defer w.watcher.Close()

	if err := w.addWatchRecursive(w.config.ContentDir); err != nil {
		return fmt.Errorf("failed to add content directory to watcher: %w", err)
	}

	if err := w.addWatchRecursive(w.config.TemplateDir); err != nil {
		return fmt.Errorf("failed to add template directory to watcher: %w", err)
	}

	if _, err := os.Stat(w.config.AssetsDir); err == nil {
		if err := w.addWatchRecursive(w.config.AssetsDir); err != nil {
			return fmt.Errorf("failed to add assets directory to watcher: %w", err)
		}
	}

	log.Printf("Watching %s, %s, and %s for changes...", w.config.ContentDir, w.config.TemplateDir, w.config.AssetsDir)

	for {
		select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return nil
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						w.addWatchRecursive(event.Name)
					}
				}

				isAsset := strings.HasPrefix(event.Name, w.config.AssetsDir)
				if isAsset {
					if event.Op&fsnotify.Remove == fsnotify.Remove {
						continue
					}
					log.Printf("Asset changed: %s, TO BE IMPLEMENTED!", event.Name)
					continue
				}

				if strings.HasPrefix(event.Name, w.config.TemplateDir) {
					if event.Op&fsnotify.Remove == fsnotify.Remove {
						continue
					}
					log.Printf("Template changed: %s, TO BE IMPLEMENTED...", event.Name)
					continue
				}

				if (!parser.IsContentFile(event.Name)) {
					continue
				}

				fmt.Print(event.Op)

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Printf("File removed: %s", event.Name)
					w.generator.HandleContentChange(generator.ContentChange{
						Op:     generator.OpDelete,
						Path: event.Name,
					})
					continue
				}


				if event.Op&fsnotify.Rename == fsnotify.Rename {
					if _, err := os.Stat(event.Name); os.IsNotExist(err) {
						log.Printf("File deleted: %s", event.Name)
						w.generator.HandleContentChange(generator.ContentChange{
							Op:   generator.OpDelete,
							Path: event.Name,
						})
					} else {
						log.Printf("File renamed from: %s", event.Name)
						w.generator.HandleContentChange(generator.ContentChange{
							Op:     generator.OpRename,
							Path: event.Name,
						})
					}
					continue
				}

				// Debounce file changes - FIND A BETTER WAY TO DO THIS!
				fileName := event.Name
				w.timerMu.Lock()
				if timer, exists := w.timers[fileName]; exists {
					timer.Stop()
				}

				w.timers[fileName] = time.AfterFunc(100*time.Millisecond, func() {
					log.Printf("File changed, rebuilding site: %s", fileName)
					w.generator.HandleContentChange(generator.ContentChange{
						Op:    generator.OpCreate,
						Path: event.Name,
					})
					delete(w.timers, fileName)
				})
				w.timerMu.Unlock()

				case err, ok := <-w.watcher.Errors:
					if !ok {
						return nil
					}
					log.Printf("Watcher error: %v", err)
		}
	}
}

func (w *Watcher) addWatchRecursive(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return w.watcher.Add(path)
		}
		return nil
	})
}
