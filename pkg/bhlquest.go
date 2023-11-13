package bhlquest

import (
	"fmt"
	"log/slog"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/bhln"
	"github.com/gnames/bhlquest/pkg/ent/embed"
	"github.com/gnames/bhlquest/pkg/ent/llmutil"
	"github.com/gnames/bhlquest/pkg/ent/storage"
)

type Components struct {
	BHLNames bhln.BHLN
	Embed    embed.Embed
	LlmUtil  llmutil.LlmUtil
	Storage  storage.Storage
}

type bhlquest struct {
	cfg  config.Config
	bhln bhln.BHLN
	emb  embed.Embed
	llm  llmutil.LlmUtil
	stg  storage.Storage
}

func New(
	cfg config.Config,
	cmp Components,
) BHLQuest {
	res := bhlquest{
		cfg:  cfg,
		bhln: cmp.BHLNames,
		emb:  cmp.Embed,
		llm:  cmp.LlmUtil,
		stg:  cmp.Storage,
	}

	return res
}

func (bq bhlquest) Init() error {
	slog.Info("Start initial data process")
	slog.Info("Collect IDs of relevant BHL items.")
	ids, err := bq.bhln.ItemIds(0, 0, nil)
	if err != nil {
		return err
	}

	slog.Info("Initiate BHLquest database.")
	err = bq.emb.Init()
	if err != nil {
		return err
	}

	slog.Info("Find Items' texts and prepare them for AI.")
	err = bq.emb.Populate(ids[0:10])
	if err != nil {
		return err
	}

	slog.Info("Initial processing finished without errors.")
	return nil
}

// GetVersion provides version information of the app.
func GetVersion() string {
	version := fmt.Sprintf("Version: %s\nBuild:   %s", Version, Build)
	return version
}
