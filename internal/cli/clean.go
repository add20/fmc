package cli

import (
	"fmt"
	"os"

	"github.com/add20/fmc/internal/config"
	"github.com/add20/fmc/internal/fmcerr"
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
				return &fmcerr.FMCError{Code: fmcerr.ErrWriteFile, Message: "failed to remove dist dir", Cause: err}
			}
			fmt.Println("clean completed.")
			return nil
		},
	}
}
