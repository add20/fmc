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

type fileInfo struct {
	srcPath string // contentsDir からの相対パス
	absPath string
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
	files, err := walkFiles(cfg.Contents.Dir)
	if err != nil {
		return err
	}
	if err := checkDuplicateSlugs(files); err != nil {
		return err
	}
	var entries []IndexEntry
	for _, f := range files {
		entry, err := compileFile(f, cfg.Output.Dir)
		if err != nil {
			return err
		}
		entries = append(entries, entry)
	}
	return writeIndex(entries, cfg.Output.Dir)
}

func walkFiles(dir string) ([]fileInfo, error) {
	var files []fileInfo
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		files = append(files, fileInfo{srcPath: rel, absPath: path})
		return nil
	})
	if err != nil {
		return nil, &fmcerr.FMCError{Code: fmcerr.ErrWriteFile, Message: "failed to walk contents dir", Cause: err}
	}
	return files, nil
}

func checkDuplicateSlugs(files []fileInfo) error {
	slugToFiles := map[string][]string{}
	for _, f := range files {
		s := Slug(f.srcPath)
		slugToFiles[s] = append(slugToFiles[s], f.absPath)
	}
	for s, paths := range slugToFiles {
		if len(paths) > 1 {
			sort.Strings(paths)
			list := strings.Join(paths, "\n- ")
			return &fmcerr.FMCError{
				Code:    fmcerr.ErrDuplicateSlug,
				Message: fmt.Sprintf("duplicate slug detected\n\nslug: %s\n\nfiles:\n- %s", s, list),
			}
		}
	}
	return nil
}

func compileFile(f fileInfo, outputDir string) (IndexEntry, error) {
	data, err := os.ReadFile(f.absPath)
	if err != nil {
		return IndexEntry{}, &fmcerr.FMCError{Code: fmcerr.ErrWriteFile, Message: "failed to read file", Cause: err}
	}

	res, err := frontmatter.Parse(string(data))
	if err != nil {
		return IndexEntry{}, &fmcerr.FMCError{Code: fmcerr.ErrFrontMatterParse, Message: f.srcPath, Cause: err}
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
		return IndexEntry{}, &fmcerr.FMCError{Code: fmcerr.ErrWriteFile, Message: "mkdir failed", Cause: err}
	}

	jsonData, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return IndexEntry{}, &fmcerr.FMCError{Code: fmcerr.ErrWriteFile, Message: "marshal failed", Cause: err}
	}
	if err := os.WriteFile(outAbs, jsonData, 0644); err != nil {
		return IndexEntry{}, &fmcerr.FMCError{Code: fmcerr.ErrWriteFile, Message: "write failed", Cause: err}
	}

	var title *string
	if t, ok := res.FrontMatter["title"]; ok {
		if s, ok := t.(string); ok {
			title = &s
		}
	}
	return IndexEntry{Slug: slug, Path: outRel, Title: title}, nil
}

func writeIndex(entries []IndexEntry, outputDir string) error {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})
	if entries == nil {
		entries = []IndexEntry{}
	}

	indexData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return &fmcerr.FMCError{Code: fmcerr.ErrWriteFile, Message: "marshal index failed", Cause: err}
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return &fmcerr.FMCError{Code: fmcerr.ErrWriteFile, Message: "mkdir dist failed", Cause: err}
	}
	if err := os.WriteFile(filepath.Join(outputDir, "index.json"), indexData, 0644); err != nil {
		return &fmcerr.FMCError{Code: fmcerr.ErrWriteFile, Message: "write index.json failed", Cause: err}
	}
	return nil
}
