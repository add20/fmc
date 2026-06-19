package watcher

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/add20/fmc/internal/compiler"
	"github.com/add20/fmc/internal/config"
	"github.com/add20/fmc/internal/fmcerr"
	"github.com/fsnotify/fsnotify"
)

func Watch(cfg config.Config) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("watcher: %w", err)
	}
	defer w.Close()

	if err := addDirs(w, cfg.Contents.Dir); err != nil {
		return err
	}

	fmt.Printf("watching %s ...\n", cfg.Contents.Dir)
	runBuild(w, cfg)

	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return nil
			}
			if event.Has(fsnotify.Create) {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					_ = w.Add(event.Name)
					continue
				}
			}
			if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) ||
				event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				runBuild(w, cfg)
			}
		case err, ok := <-w.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		}
	}
}

func runBuild(w *fsnotify.Watcher, cfg config.Config) {
	err := compiler.Build(cfg)
	if err == nil {
		fmt.Println("build completed.")
		return
	}

	var fmcErr *fmcerr.FMCError
	if errors.As(err, &fmcErr) {
		switch fmcErr.Code {
		case fmcerr.ErrDuplicateSlug:
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			_ = os.Remove(filepath.Join(cfg.Output.Dir, "index.json"))
			return
		case fmcerr.ErrFrontMatterParse:
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			// fmcErr.Message には srcPath が入っている
			jsonPath := filepath.Join(cfg.Output.Dir, fmcErr.Message+".json")
			_ = os.Remove(jsonPath)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
}

func addDirs(w *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return w.Add(path)
		}
		return nil
	})
}
