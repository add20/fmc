package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const defaultConfig = `[contents]
dir = "contents"

[output]
dir = "dist"
`

func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "初期ディレクトリ構成を生成する",
		RunE: func(cmd *cobra.Command, args []string) error {
			dirs := []string{"contents", "dist", "settings"}
			for _, d := range dirs {
				if err := mkdirIfNotExist(d); err != nil {
					return err
				}
			}
			if err := writeIfNotExist("settings/config.toml", defaultConfig); err != nil {
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
