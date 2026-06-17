package compiler

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/add20/fmc/internal/config"
	"github.com/add20/fmc/internal/fmcerr"
	"github.com/add20/fmc/internal/frontmatter"
)

// fmcerr の型をこのパッケージからも参照できるようにエイリアスを公開する
type FMCError = fmcerr.FMCError
type ErrorCode = fmcerr.ErrorCode

const (
	ErrDuplicateSlug    = fmcerr.ErrDuplicateSlug
	ErrFrontMatterParse = fmcerr.ErrFrontMatterParse
	ErrConfigLoad       = fmcerr.ErrConfigLoad
	ErrWriteFile        = fmcerr.ErrWriteFile
)

type Document struct {
	Slug        string         `json:"slug"`
	SrcPath     string         `json:"srcPath"`
	FrontMatter map[string]any `json:"frontMatter"`
	Content     string         `json:"content"`
}

type IndexEntry struct {
	Slug  string  `json:"slug"`
	Path  string  `json:"path"`
	Title *string `json:"title"`
}

func Slug(filename string) string {
	name := filepath.Base(filename)
	for {
		ext := filepath.Ext(name)
		if ext == "" {
			break
		}
		name = strings.TrimSuffix(name, ext)
	}
	return name
}

func Build(cfg config.Config) error {
	contentsDir := cfg.Contents.Dir
	outputDir := cfg.Output.Dir

	type fileInfo struct {
		srcPath string // relative to contentsDir
		absPath string
	}

	var files []fileInfo
	err := filepath.WalkDir(contentsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(contentsDir, path)
		files = append(files, fileInfo{srcPath: rel, absPath: path})
		return nil
	})
	if err != nil {
		return &FMCError{Code: ErrWriteFile, Message: "failed to walk contents dir", Cause: err}
	}

	// duplicate slug check
	slugToFiles := map[string][]string{}
	for _, f := range files {
		s := Slug(f.srcPath)
		slugToFiles[s] = append(slugToFiles[s], filepath.Join(contentsDir, f.srcPath))
	}
	for s, paths := range slugToFiles {
		if len(paths) > 1 {
			sort.Strings(paths)
			list := strings.Join(paths, "\n- ")
			return &FMCError{
				Code:    ErrDuplicateSlug,
				Message: fmt.Sprintf("duplicate slug detected\n\nslug: %s\n\nfiles:\n- %s", s, list),
			}
		}
	}

	// parse and write
	var entries []IndexEntry
	for _, f := range files {
		data, err := os.ReadFile(f.absPath)
		if err != nil {
			return &FMCError{Code: ErrWriteFile, Message: "failed to read file", Cause: err}
		}

		res, err := frontmatter.Parse(string(data))
		if err != nil {
			return &FMCError{Code: ErrFrontMatterParse, Message: f.srcPath, Cause: err}
		}

		slug := Slug(f.srcPath)
		doc := Document{
			Slug:        slug,
			SrcPath:     filepath.ToSlash(f.srcPath),
			FrontMatter: res.FrontMatter,
			Content:     res.Content,
		}

		outRel := filepath.ToSlash(f.srcPath) + ".json"
		outAbs := filepath.Join(outputDir, f.srcPath+".json")
		if err := os.MkdirAll(filepath.Dir(outAbs), 0755); err != nil {
			return &FMCError{Code: ErrWriteFile, Message: "mkdir failed", Cause: err}
		}

		jsonData, _ := json.MarshalIndent(doc, "", "  ")
		if err := os.WriteFile(outAbs, jsonData, 0644); err != nil {
			return &FMCError{Code: ErrWriteFile, Message: "write failed", Cause: err}
		}

		var title *string
		if t, ok := res.FrontMatter["title"]; ok {
			if s, ok := t.(string); ok {
				title = &s
			}
		}
		entries = append(entries, IndexEntry{Slug: slug, Path: outRel, Title: title})
	}

	// sort by srcPath
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	if entries == nil {
		entries = []IndexEntry{}
	}
	indexData, _ := json.MarshalIndent(entries, "", "  ")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return &FMCError{Code: ErrWriteFile, Message: "mkdir dist failed", Cause: err}
	}
	if err := os.WriteFile(filepath.Join(outputDir, "index.json"), indexData, 0644); err != nil {
		return &FMCError{Code: ErrWriteFile, Message: "write index.json failed", Cause: err}
	}

	return nil
}
