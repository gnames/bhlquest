package embed

import "github.com/gnames/bhlquest/pkg/ent/answer"

type Embed interface {
	Init() error
	Populate(itemIDs []uint) error
	Embed(q []string) ([][]float32, error)
	Query(emb []float32) (answer.Answer, error)
}
