package cmd

import (
	"log/slog"
	"os"

	"github.com/gnames/bhlquest/internal/gpt"
	"github.com/gnames/bhlquest/internal/io/bhlnio"
	"github.com/gnames/bhlquest/internal/io/cohere"
	"github.com/gnames/bhlquest/internal/io/gptio"
	"github.com/gnames/bhlquest/internal/io/llmutilio"
	"github.com/gnames/bhlquest/internal/io/qdrant"
	"github.com/gnames/bhlquest/internal/io/storageio"
	bhlquest "github.com/gnames/bhlquest/pkg"
	"github.com/gnames/bhlquest/pkg/config"
)

func bhlquestFactory() bhlquest.BHLQuest {
	cfg := config.New(opts...)

	gAPI := gptio.New(cfg)
	g := gpt.New(cfg, gAPI)

	bn, err := bhlnio.New(cfg)
	if err != nil {
		slog.Error("No connection to BHLNames db", "error", err)
		os.Exit(1)
	}

	llm, err := llmutilio.New(cfg)
	if err != nil {
		slog.Error("No connection to llmutil", "error", err)
		os.Exit(1)
	}

	stg := storageio.New(cfg)
	emb, err := qdrant.New(cfg, stg, llm)
	if err != nil {
		slog.Error("No connection to embedding db", "error", err)
		os.Exit(1)
	}

	cp := bhlquest.Components{
		BHLN:     bn,
		Embed:    emb,
		Reranker: cohere.New(cfg),
		GPT:      g,
	}

	bq := bhlquest.New(cfg, cp)
	return bq
}
