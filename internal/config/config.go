package config

import (
	"os"

	"github.com/add20/fmc/internal/fmcerr"
	"github.com/pelletier/go-toml/v2"
)

const DefaultConfigPath = "settings/config.toml"

type Config struct {
	Contents struct {
		Dir string `toml:"dir"`
	} `toml:"contents"`

	Output struct {
		Dir string `toml:"dir"`
	} `toml:"output"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, &fmcerr.FMCError{Code: fmcerr.ErrConfigLoad, Message: path, Cause: err}
	}
	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return Config{}, &fmcerr.FMCError{Code: fmcerr.ErrConfigLoad, Message: "failed to parse config", Cause: err}
	}
	return cfg, nil
}
