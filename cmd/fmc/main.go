package main

import (
	"os"

	"github.com/add20/fmc/internal/cli"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "fmc",
		Short: "FrontMatter Compiler",
	}
	root.AddCommand(
		cli.NewBuildCmd(),
		cli.NewInitCmd(),
		cli.NewCleanCmd(),
		cli.NewWatchCmd(),
	)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
