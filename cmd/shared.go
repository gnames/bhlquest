package cmd

import (
	"log/slog"
	"os"

	"github.com/gnames/bhlquest/internal/io/bhlnio"
	"github.com/gnames/bhlquest/internal/io/embedio"
	"github.com/gnames/bhlquest/internal/io/llmutilio"
	"github.com/gnames/bhlquest/internal/io/storageio"
	bhlquest "github.com/gnames/bhlquest/pkg"
	"github.com/gnames/bhlquest/pkg/config"
)

func bhlquestFactory() bhlquest.BHLQuest {
	cfg := config.New(opts...)

	bn, err := bhlnio.New(cfg)
	if err != nil {
		slog.Error("No connection to BHLNames db", "error", err)
		os.Exit(1)
	}

	stg := storageio.New(cfg)

	llm, err := llmutilio.New(cfg)
	if err != nil {
		slog.Error("No connection to llmutil", "error", err)
		os.Exit(1)
	}
	emb, err := embedio.New(cfg, stg, llm)
	if err != nil {
		slog.Error("No connection to embedding db", "error", err)
		os.Exit(1)
	}

	cp := bhlquest.Components{
		BHLNames: bn,
		Embed:    emb,
	}

	bq := bhlquest.New(cfg, cp)
	return bq

}
