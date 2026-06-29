package compiler_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/add20/fmc/internal/compiler"
	"github.com/add20/fmc/internal/config"
	"github.com/add20/fmc/internal/fmcerr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testConfig(t *testing.T) (config.Config, string) {
	t.Helper()
	distDir := t.TempDir()
	var cfg config.Config
	cfg.Contents.Dir = "../../testdata/contents"
	cfg.Output.Dir = distDir
	return cfg, distDir
}

func TestSlug(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"README.md", "README"},
		{"README.txt", "README"},
		{"README.md.txt", "README"},
		{"archive.tar.gz", "archive"},
		{"memo", "memo"},
		{"2026/06/README.md", "2026/06/README"},
		{"2026/06/archive.tar.gz", "2026/06/archive"},
		{"a/b/c/memo", "a/b/c/memo"},
	}
	for _, c := range cases {
		assert.Equal(t, c.out, compiler.Slug(c.in), c.in)
	}
}

func TestBuild(t *testing.T) {
	cfg, distDir := testConfig(t)
	err := compiler.Build(cfg)
	require.NoError(t, err)

	// index.json が生成されているか
	indexPath := filepath.Join(distDir, "index.json")
	data, err := os.ReadFile(indexPath)
	require.NoError(t, err)

	var entries []compiler.IndexEntry
	require.NoError(t, json.Unmarshal(data, &entries))
	assert.GreaterOrEqual(t, len(entries), 1)

	// title が null になるエントリの確認
	var foundNoTitle bool
	for _, e := range entries {
		if e.Slug == "notitle" {
			assert.Nil(t, e.Title)
			assert.False(t, e.Draft)
			foundNoTitle = true
		}
	}
	assert.True(t, foundNoTitle, "notitle エントリが見つからない")

	// draft: true のエントリの確認
	var foundDraft bool
	for _, e := range entries {
		if e.Slug == "draft-post" {
			assert.True(t, e.Draft)
			foundDraft = true
		}
	}
	assert.True(t, foundDraft, "draft-post エントリが見つからない")
}

func TestBuildDuplicateSlug(t *testing.T) {
	// 重複 slug を引き起こす一時ディレクトリ
	contentsDir := t.TempDir()
	os.WriteFile(filepath.Join(contentsDir, "README.md"), []byte("---\ntitle: A\n---\n"), 0644)
	os.WriteFile(filepath.Join(contentsDir, "README.txt"), []byte("---\ntitle: B\n---\n"), 0644)

	var cfg config.Config
	cfg.Contents.Dir = contentsDir
	cfg.Output.Dir = t.TempDir()

	err := compiler.Build(cfg)
	require.Error(t, err)
	fmcErr, ok := err.(*fmcerr.FMCError)
	require.True(t, ok)
	assert.Equal(t, fmcerr.ErrDuplicateSlug, fmcErr.Code)
}

func TestBuildOutputJSON(t *testing.T) {
	cfg, distDir := testConfig(t)
	require.NoError(t, compiler.Build(cfg))

	// README.md.json の内容を検証
	jsonPath := filepath.Join(distDir, "2026/06/README.md.json")
	data, err := os.ReadFile(jsonPath)
	require.NoError(t, err)

	var doc compiler.Document
	require.NoError(t, json.Unmarshal(data, &doc))
	assert.Equal(t, "2026/06/README", doc.Slug)
	assert.Equal(t, "2026/06/README.md", doc.SrcPath)
	assert.Equal(t, "最初に読むファイル", doc.FrontMatter["title"])
}
