package embed

import (
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/output"
)

type Embed interface {
	Init() error
	SetItemsNum(int)
	Populate(itemIDs []uint) error
	Embed(q []string) ([][]float32, error)
	CrossEmbed(ss [][]string) ([]float64, error)
	Query(emb []float32) (output.Answer, error)
	SetConfig(config.Config) Embed
}
