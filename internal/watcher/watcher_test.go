package watcher

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/add20/fmc/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeConfig(t *testing.T, contentsDir string) config.Config {
	t.Helper()
	var cfg config.Config
	cfg.Contents.Dir = contentsDir
	cfg.Output.Dir = t.TempDir()
	return cfg
}

func TestRunBuildSuccess(t *testing.T) {
	contentsDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(contentsDir, "hello.md"), []byte("---\ntitle: Hello\n---\nbody"), 0644))

	cfg := makeConfig(t, contentsDir)
	runBuild(cfg)

	_, err := os.Stat(filepath.Join(cfg.Output.Dir, "index.json"))
	assert.NoError(t, err, "index.json が生成されていること")
}

func TestRunBuildDuplicateSlugRemovesIndex(t *testing.T) {
	contentsDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(contentsDir, "README.md"), []byte("---\ntitle: A\n---\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(contentsDir, "README.txt"), []byte("---\ntitle: B\n---\n"), 0644))

	cfg := makeConfig(t, contentsDir)

	// 事前に index.json を置いておく（直前ビルドの残骸を再現）
	require.NoError(t, os.MkdirAll(cfg.Output.Dir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(cfg.Output.Dir, "index.json"), []byte("[]"), 0644))

	runBuild(cfg)

	_, err := os.Stat(filepath.Join(cfg.Output.Dir, "index.json"))
	assert.True(t, os.IsNotExist(err), "slug 重複時は index.json が削除されること")
}

func TestRunBuildFrontMatterParseErrorRemovesJSON(t *testing.T) {
	contentsDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(contentsDir, "bad.md"), []byte("---\n: invalid\n---\n"), 0644))

	cfg := makeConfig(t, contentsDir)

	// 事前に対応 JSON を置いておく（直前ビルドの残骸を再現）
	require.NoError(t, os.MkdirAll(cfg.Output.Dir, 0755))
	staleJSON := filepath.Join(cfg.Output.Dir, "bad.md.json")
	require.NoError(t, os.WriteFile(staleJSON, []byte("{}"), 0644))

	runBuild(cfg)

	_, err := os.Stat(staleJSON)
	assert.True(t, os.IsNotExist(err), "パースエラー時は該当 JSON が削除されること")
}
