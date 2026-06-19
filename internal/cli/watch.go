package cli

import (
	"github.com/add20/fmc/internal/config"
	"github.com/add20/fmc/internal/watcher"
	"github.com/spf13/cobra"
)

func NewWatchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "watch",
		Short: "contents ディレクトリを監視し、変更時に自動ビルドする",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultConfigPath)
			if err != nil {
				return err
			}
			return watcher.Watch(cfg)
		},
	}
}
