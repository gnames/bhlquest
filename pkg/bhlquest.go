package bhlquest

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gnames/bhlquest/internal/bhln"
	"github.com/gnames/bhlquest/internal/embed"
	"github.com/gnames/bhlquest/internal/ent/ref"
	"github.com/gnames/bhlquest/internal/gpt"
	"github.com/gnames/bhlquest/internal/rerank"
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/output"
	"github.com/gnames/gnlib"
	"github.com/gnames/gnlib/ent/gnvers"
)

type Components struct {
	bhln.BHLN
	embed.Embed
	rerank.Reranker
	gpt.GPT
}

type bhlquest struct {
	cfg  config.Config
	bhln bhln.BHLN
	emb  embed.Embed
	rnk  rerank.Reranker
	gpt  gpt.GPT
}

// New creates a new instance of BHLQuest.
func New(
	cfg config.Config,
	cmp Components,
) BHLQuest {
	res := bhlquest{
		cfg:  cfg,
		bhln: cmp.BHLN,
		emb:  cmp.Embed,
		gpt:  cmp.GPT,
		rnk:  cmp.Reranker,
	}

	return res
}

// Init creates a vector data-store from BHL Items.
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

func (bq bhlquest) Ask(q string) (output.Answer, error) {
	start := time.Now()
	var res output.Answer
	var results []*output.Result
	emb, err := bq.emb.Embed([]string{q})
	if err != nil {
		return res, err
	}
	if len(emb) < 1 {
		err = errors.New("embedding of the question failed")
		return res, err
	}

	res, err = bq.emb.Query(emb[0])
	if err != nil {
		return res, err
	}

	results, err = bq.addReferences(res.Results)
	if err != nil {
		return res, err
	}

	results, err = bq.addPageNums(results)
	if err != nil {
		return res, err
	}

	results, err = bq.rnk.Rerank(q, results)
	if err != nil {
		return res, err
	}
	res.Results = results

	if len(res.Results) > bq.cfg.MaxResultsNum {
		res.Results = res.Results[:bq.cfg.MaxResultsNum]
	}
	duration := time.Since(start).Seconds()
	res.Question = q
	res.QueryTime = duration
	res.Version = GetVersion().Version
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

// GetConfig provides configuration of the app.
func (bq bhlquest) GetConfig() config.Config {
	return bq.cfg
}

// SetConfig updates configuration of the app.
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

func (bq bhlquest) addPageNums(
	results []*output.Result,
) ([]*output.Result, error) {
	rs := results
	for i := range rs {
		pages, err := bq.bhln.PageNums(rs[i].ItemID)
		if err != nil {
			return nil, err
		}

		for j := range rs[i].Pages {
			if pageNum, ok := pages[rs[i].Pages[j].ID]; ok {
				rs[i].Pages[j].PageNum = pageNum
			}
		}

	}
	return rs, nil
}

// addReferences adds references to the results, it also tries to remove duplicates
// using ref.Reference.Fingerprint.
func (bq bhlquest) addReferences(
	results []*output.Result,
) ([]*output.Result, error) {
	var res []*output.Result
	ids := gnlib.Map(results, func(r *output.Result) int {
		return int(r.PageID)
	})
	refs, err := bq.bhln.References(ids)
	if err != nil {
		return res, err
	}
	rs := results
	dupl := make(map[string]struct{})
	var ref ref.Reference
	var ok bool
	for i := range rs {
		if ref, ok = refs[int(rs[i].PageID)]; !ok {
			res = append(res, rs[i])
			continue
		}
		if _, ok = dupl[ref.Fingerprint]; ok {
			continue
		}
		dupl[ref.Fingerprint] = struct{}{}
		rs[i].Reference = ref.String()
		if ref.TitleDOI != "" {
			rs[i].OutlinkTitleDOI = "https://doi.org/" + ref.TitleDOI

		}
		if ref.TitleLang != "" {
			rs[i].Language = strings.ToLower(ref.TitleLang)
		}
		res = append(res, rs[i])
	}
	return res, nil
}
