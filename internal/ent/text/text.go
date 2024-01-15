package text

import (
	"github.com/gnames/bhlquest/internal/storage"
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/gnbhl/itembhl"
	"github.com/gnames/gnlib"
	"github.com/gnames/gnlib/ent/gnml"
)

type text struct {
	cfg config.Config
	stg storage.Storage
}

func New(cfg config.Config, stg storage.Storage) Text {
	res := &text{
		cfg: cfg,
		stg: stg,
	}

	return res
}

func (t *text) ItemToChunks(itemID uint) ([]Chunk, error) {
	itm, err := itembhl.New(t.cfg.BHLDir, itemID)
	if err != nil {
		return nil, err
	}
	txt := itm.Text()

	txts := gnml.SplitText(txt, 1500, 150)
	chunks := gnlib.Map(txts, func(tp gnml.TextPart) Chunk {
		start := tp.StartOffset
		res := Chunk{
			ItemID: itemID,
			Text:   tp.Content,
			Start:  uint(start),
			Length: uint(tp.Length),
		}
		return res
	})

	return chunks, nil
}

func (t *text) ItemByID(itemID uint) (itembhl.ItemBHL, error) {
	return itembhl.New(t.cfg.BHLDir, itemID)
}
