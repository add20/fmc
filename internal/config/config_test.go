package config_test

import (
	"errors"
	"testing"

	"github.com/add20/fmc/internal/config"
	"github.com/add20/fmc/internal/fmcerr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	cfg, err := config.Load("../../testdata/configs/default.toml")
	require.NoError(t, err)
	assert.Equal(t, "testdata/contents", cfg.Contents.Dir)
	assert.Equal(t, "testdata/dist", cfg.Output.Dir)
}

func TestLoadIndexFields(t *testing.T) {
	cfg, err := config.Load("../../testdata/configs/with_index_fields.toml")
	require.NoError(t, err)
	assert.Equal(t, []string{"category", "tags"}, cfg.Index.Fields)
}

func TestLoadNotFound(t *testing.T) {
	_, err := config.Load("nonexistent.toml")
	require.Error(t, err)
	var fmcErr *fmcerr.FMCError
	require.True(t, errors.As(err, &fmcErr))
	assert.Equal(t, fmcerr.ErrConfigLoad, fmcErr.Code)
}
