package embed

import (
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/answer"
)

type Embed interface {
	Init() error
	LastItemID() uint
	SetItemsNum(int)
	Populate(itemIDs []uint) error
	Embed(q []string) ([][]float32, error)
	CrossEmbed(ss [][]string) ([]float64, error)
	Query(emb []float32) (answer.Answer, error)
	SetConfig(config.Config) Embed
}
