package cli

import (
	"fmt"
	"os"

	"github.com/add20/fmc/internal/config"
	"github.com/spf13/cobra"
)

var defaultConfigContent = fmt.Sprintf("[contents]\ndir = %q\n\n[output]\ndir = %q\n", config.DefaultContentsDir, config.DefaultOutputDir)

func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "初期ディレクトリ構成を生成する",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, d := range []string{config.DefaultContentsDir, config.DefaultOutputDir, config.DefaultSettingsDir} {
				if err := mkdirIfNotExist(d); err != nil {
					return err
				}
			}
			if err := writeIfNotExist(config.DefaultConfigPath, defaultConfigContent); err != nil {
				return err
			}
			fmt.Println("init completed.")
			return nil
		},
	}
}

func mkdirIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func writeIfNotExist(path, content string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.WriteFile(path, []byte(content), 0644)
	}
	return nil
}
