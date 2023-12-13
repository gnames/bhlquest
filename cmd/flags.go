package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	bhlquest "github.com/gnames/bhlquest/pkg"
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/gnlib"
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

func rebuildDbFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("rebuild-db")
	if b {
		opts = append(opts, config.OptWithRebuildDb(true))
	}
}

func noConfirmFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("yes-to-confirmations")
	if b {
		opts = append(opts, config.OptWithoutConfirm(true))
	}
}

func portFlag(cmd *cobra.Command) {
	i, _ := cmd.Flags().GetInt("port")
	if i > 0 {
		opts = append(opts, config.OptPort(i))
	}
}

func maxResultsFlag(cmd *cobra.Command) {
	i, _ := cmd.Flags().GetInt("max-results")
	if i > 0 {
		opts = append(opts, config.OptMaxResultsNum(i))
	}
}

func scoreThresholdFlag(cmd *cobra.Command) {
	f, _ := cmd.Flags().GetFloat64("score-threshold")
	if f > 1 || f < 0 {
		slog.Error("The score-threshold should be betwen 0 and 1")
		slog.Error("Score 1 is perfect, 0 is no match.")
		f = 0
	}
	if f > 0 {
		opts = append(opts, config.OptScoreThreshold(1-f))
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

func classesFlag(cmd *cobra.Command) {
	s, _ := cmd.Flags().GetString("classes")
	if s != "" {
		el := strings.Split(s, ",")
		classes := gnlib.Map(el, func(s string) string {
			return strings.TrimSpace(s)
		})
		opts = append(opts, config.OptInitClasses(classes))
	}
}

func taxonsFlag(cmd *cobra.Command) {
	s, _ := cmd.Flags().GetString("taxons")
	if s != "" {
		el := strings.Split(s, ",")
		taxa := gnlib.Map(el, func(s string) string {
			return strings.TrimSpace(s)
		})
		opts = append(opts, config.OptInitTaxa(taxa))
	}
}
