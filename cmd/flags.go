package cmd

import (
	"fmt"
	"log/slog"
	"os"

	bhlquest "github.com/gnames/bhlquest/pkg"
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/spf13/cobra"
)

type flagFunc func(cmd *cobra.Command)

func debugFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("debug")
	if b {
		opts := &slog.HandlerOptions{Level: slog.LevelDebug}
		handle := slog.NewJSONHandler(os.Stderr, opts)
		logger := slog.New(handle)
		slog.SetDefault(logger)
	}
}

func portFlag(cmd *cobra.Command) {
	i, _ := cmd.Flags().GetInt("port")
	if i > 0 {
		opts = append(opts, config.OptPort(i))
	}
}

func versionFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("version")
	if b {
		version := bhlquest.GetVersion()
		fmt.Printf(
			"\nVersion: %s\nBuild:   %s\n\n",
			version.Version,
			version.Build,
		)
		os.Exit(0)
	}
}
