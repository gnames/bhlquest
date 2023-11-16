package bhlquest

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/answer"
	"github.com/gnames/bhlquest/pkg/ent/bhln"
	"github.com/gnames/bhlquest/pkg/ent/embed"
	"github.com/gnames/gnlib/ent/gnvers"
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
	err = bq.emb.Populate(ids)
	if err != nil {
		return err
	}

	slog.Info("Initial processing finished without errors.")
	return nil
}

func (bq bhlquest) Ask(q string) (answer.Answer, error) {
	start := time.Now()
	var res answer.Answer
	emb, err := bq.emb.Embed([]string{q})
	if err != nil {
		return res, err
	}
	if len(emb) < 1 {
		err := errors.New("embedding of the question failed")
		return res, err
	}
	res, err = bq.emb.Query(emb[0])
	if err != nil {
		return res, err
	}
	duration := time.Since(start).Seconds()
	res.Meta.Question = q
	res.Meta.QueryTime = duration
	return res, nil
}

func (bq bhlquest) GetConfig() config.Config {
	return bq.cfg
}

// GetVersion provides version information of the app.
func GetVersion() gnvers.Version {
	version := gnvers.Version{
		Version: Version,
		Build:   Build,
	}
	return version
}
