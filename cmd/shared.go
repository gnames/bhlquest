package cmd

import (
	"fmt"
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
		msg := fmt.Sprintf("Init of BHLNames failed: %s", err)
		slog.Error(msg)
		os.Exit(1)
	}

	stg := storageio.New(cfg)

	llm, err := llmutilio.New(cfg)
	if err != nil {
		msg := fmt.Sprintf("Init of BHLNames failed: %s", err)
		slog.Error(msg)
		os.Exit(1)
	}
	emb, err := embedio.New(cfg, stg, llm)
	if err != nil {
		msg := fmt.Sprintf("Init of BHLNames failed: %s", err)
		slog.Error(msg)
		os.Exit(1)
	}

	cp := bhlquest.Components{
		BHLNames: bn,
		Embed:    emb,
	}

	bq := bhlquest.New(cfg, cp)
	return bq

}
