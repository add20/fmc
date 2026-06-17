package config_test

import (
	"testing"

	"github.com/add20/fmc/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	cfg, err := config.Load("../../testdata/configs/default.toml")
	require.NoError(t, err)
	assert.Equal(t, "testdata/contents", cfg.Contents.Dir)
	assert.Equal(t, "testdata/dist", cfg.Output.Dir)
}

func TestLoadNotFound(t *testing.T) {
	_, err := config.Load("nonexistent.toml")
	assert.Error(t, err)
}
