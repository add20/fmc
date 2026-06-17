package cli

import (
	"fmt"
	"os"

	"github.com/add20/fmc/internal/config"
	"github.com/spf13/cobra"
)

func NewCleanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "dist ディレクトリを削除する",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultConfigPath)
			if err != nil {
				return err
			}
			if err := os.RemoveAll(cfg.Output.Dir); err != nil {
				return err
			}
			fmt.Println("clean completed.")
			return nil
		},
	}
}
