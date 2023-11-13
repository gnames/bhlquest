package bhlquest

import (
	"fmt"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/bhln"
	"github.com/gnames/bhlquest/pkg/ent/embed"
)

type bhlquest struct {
	cfg  config.Config
	bhln bhln.BHLN
	emb  embed.Embed
}

func New(
	cfg config.Config,
	bhln bhln.BHLN,
	emb embed.Embed,
) BHLQuest {
	res := bhlquest{
		cfg:  cfg,
		bhln: bhln,
		emb:  emb,
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

	_ = ids
	return nil
}

// GetVersion provides version information of the app.
func GetVersion() string {
	version := fmt.Sprintf("Version: %s\nBuild:   %s", Version, Build)
	return version
}
