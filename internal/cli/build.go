package cli

import (
	"fmt"

	"github.com/add20/fmc/internal/compiler"
	"github.com/add20/fmc/internal/config"
	"github.com/spf13/cobra"
)

func NewBuildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Frontmatter ファイルを JSON へコンパイルする",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultConfigPath)
			if err != nil {
				return err
			}
			if err := compiler.Build(cfg); err != nil {
				return fmt.Errorf("%w", err)
			}
			fmt.Println("build completed.")
			return nil
		},
	}
}
