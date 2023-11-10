package cmd

import (
	"fmt"
	"log/slog"
	"os"

	bhlquest "github.com/gnames/bhlquest/pkg"
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

func versionFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("version")
	if b {
		version := bhlquest.GetVersion()
		fmt.Print(version)
		os.Exit(0)
	}
}
