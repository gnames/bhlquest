package bhlquest

import (
	"fmt"
	"log/slog"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/bhln"
	"github.com/gnames/bhlquest/pkg/ent/embed"
)

type Components struct {
	BHLNames bhln.BHLN
	Embed    embed.Embed
}

type bhlquest struct {
	cfg  config.Config
	bhln bhln.BHLN
	emb  embed.Embed
}

func New(
	cfg config.Config,
	cmp Components,
) BHLQuest {
	res := bhlquest{
		cfg:  cfg,
		bhln: cmp.BHLNames,
		emb:  cmp.Embed,
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

	msg := fmt.Sprintf("Found %d items.", len(ids))
	slog.Info(msg)

	slog.Info("Initiate BHLquest database.")
	err = bq.emb.Init()
	if err != nil {
		return err
	}

	slog.Info("Find Items' texts and prepare them for AI.")
	err = bq.emb.Populate(ids[0:1000])
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
