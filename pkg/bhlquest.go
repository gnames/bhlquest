package bhlquest

import (
	"cmp"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/answer"
	"github.com/gnames/bhlquest/pkg/ent/bhln"
	"github.com/gnames/bhlquest/pkg/ent/embed"
	"github.com/gnames/bhlquest/pkg/ent/gpt"
	"github.com/gnames/gnlib"
	"github.com/gnames/gnlib/ent/gnvers"
)

type Components struct {
	BHLNames bhln.BHLN
	Embed    embed.Embed
	GPT      gpt.GPT
}

type bhlquest struct {
	cfg  config.Config
	bhln bhln.BHLN
	emb  embed.Embed
	gpt  gpt.GPT
}

func New(
	cfg config.Config,
	cmp Components,
) BHLQuest {
	res := bhlquest{
		cfg:  cfg,
		bhln: cmp.BHLNames,
		emb:  cmp.Embed,
		gpt:  cmp.GPT,
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

	itemsNum := len(ids)
	bq.emb.SetItemsNum(itemsNum)
	slog.Info("Aquired BHL items", "items_num", itemsNum)

	if !bq.cfg.WithoutConfirm {
		fmt.Printf(
			"\nReady to process %d items. It might take a long time.\n",
			len(ids),
		)
		fmt.Println("Do you want to proceed? (y/N)")
		var confirm string
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) != "y" {
			os.Exit(0)
		}
	}

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
	if bq.cfg.WithCrossEmbed {
		cePairs := gnlib.Map(
			res.Results, func(r answer.Result) []string {
				return []string{res.Question, r.TextExt}
			},
		)
		ceScore, err := bq.emb.CrossEmbed(cePairs)
		if err != nil {
			slog.Warn("Cross-embedding failed.")
		}
		for i := range res.Results {
			res.Results[i].CrossScore = ceScore[i]
		}
		slices.SortStableFunc(
			res.Results,
			func(a, b answer.Result) int {
				res := cmp.Compare(b.CrossScore, a.CrossScore)
				if res != 0 {
					return res
				}
				return cmp.Compare(b.Score, a.Score)
			},
		)
	}
	if len(res.Results) > bq.cfg.MaxResultsNum {
		res.Results = res.Results[:bq.cfg.MaxResultsNum]
	}
	duration := time.Since(start).Seconds()
	res.Meta.Question = q
	res.Meta.QueryTime = duration
	res.Meta.Version = GetVersion().Version
	if bq.cfg.WithSummary {
		sum, err := bq.gpt.Summary(res)
		if err == nil {
			res.Summary = sum
		} else {
			slog.Warn("Summary failed: %s", err)
		}
	}
	return res, nil
}

func (bq bhlquest) GetConfig() config.Config {
	return bq.cfg
}

func (bq bhlquest) SetConfig(cfg config.Config) BHLQuest {
	bq.cfg = cfg
	bq.emb = bq.emb.SetConfig(cfg)
	return bq
}

// GetVersion provides version information of the app.
func GetVersion() gnvers.Version {
	version := gnvers.Version{
		Version: Version,
		Build:   Build,
	}
	return version
}
