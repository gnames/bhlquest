package bhlquest

import (
	"fmt"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/bhln"
	"github.com/gnames/bhlquest/pkg/ent/embed"
	"github.com/gnames/bhlquest/pkg/ent/llmutil"
)

type bhlquest struct {
	cfg  config.Config
	bhln bhln.BHLN
	emb  embed.Embed
	llm  llmutil.LlmUtil
}

func New(
	cfg config.Config,
	bhln bhln.BHLN,
	emb embed.Embed,
	llm llmutil.LlmUtil,
) BHLQuest {
	res := bhlquest{
		cfg:  cfg,
		bhln: bhln,
		emb:  emb,
		llm:  llm,
	}

	return res
}

func (bq bhlquest) Init() error {

	ids, err := bq.bhln.ItemIds(0, 0, nil)
	if err != nil {
		return err
	}

	err = bq.emb.Init()
	if err != nil {
		return err
	}

	err = bq.emb.Populate(ids)

	_ = ids
	return nil
}

// GetVersion provides version information of the app.
func GetVersion() string {
	version := fmt.Sprintf("Version: %s\nBuild:   %s", Version, Build)
	return version
}
