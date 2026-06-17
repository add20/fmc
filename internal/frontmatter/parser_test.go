package frontmatter_test

import (
	"testing"

	"github.com/add20/fmc/internal/frontmatter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseYAML(t *testing.T) {
	src := "---\ntitle: Hello\ntags:\n  - go\n---\n本文\n"
	res, err := frontmatter.Parse(src)
	require.NoError(t, err)
	assert.Equal(t, "Hello", res.FrontMatter["title"])
	assert.Equal(t, "本文\n", res.Content)
}

func TestParseTOML(t *testing.T) {
	src := "+++\ntitle = \"Hello\"\n+++\n本文\n"
	res, err := frontmatter.Parse(src)
	require.NoError(t, err)
	assert.Equal(t, "Hello", res.FrontMatter["title"])
	assert.Equal(t, "本文\n", res.Content)
}

func TestParseNoFrontmatter(t *testing.T) {
	src := "ただの本文"
	res, err := frontmatter.Parse(src)
	require.NoError(t, err)
	assert.Empty(t, res.FrontMatter)
	assert.Equal(t, src, res.Content)
}

func TestParseYAMLNoTitle(t *testing.T) {
	src := "---\ntags:\n  - test\n---\n本文\n"
	res, err := frontmatter.Parse(src)
	require.NoError(t, err)
	assert.Nil(t, res.FrontMatter["title"])
}

func TestParseYAMLInvalid(t *testing.T) {
	src := "---\n: invalid: yaml:\n---\n本文\n"
	_, err := frontmatter.Parse(src)
	assert.Error(t, err)
}
