package text

import (
	"strings"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/storage"
	"github.com/gnames/gnlib"
	"github.com/gnames/gnlib/ent/gnml"
)

type ItemText struct {
	Path    string
	TitleID uint
	PageIDs []uint
}

type Chunk struct {
	ID        uint
	ItemID    uint
	PageIDs   []uint
	Start     uint
	End       uint
	Length    uint
	Text      string
	Embedding []float32
	Score     float64
}

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

func (t *text) TextToChunks(itemID uint) ([]Chunk, error) {
	var res []Chunk
	pages, err := t.stg.Pages(itemID)
	if err != nil {
		return res, err
	}

	return t.pagesToChunks(pages), nil
}

func (t *text) ChunkText(cnk Chunk) (string, string) {
	pages, err := t.stg.Pages(cnk.ItemID)
	if err != nil {
		return "", ""
	}
	txt := combinePages(pages)
	offset := int(cnk.Start)
	offsetExt := offset - 500
	if offsetExt < 0 {
		offsetExt = 0
	}

	if len(txt) <= offset {
		return "", ""
	}
	res := txt[offset:]
	resExt := txt[offsetExt:]

	if len(resExt) >= 2000 {
		resExt = resExt[:2000]
	}

	if len(res) < 1000 {
		return res, resExt
	}
	return res[:1000], resExt
}

func (t *text) pagesToChunks(pages []storage.Page) []Chunk {
	txt := combinePages(pages)
	txts := gnml.SplitText(txt, 1500, 150)
	chunks := gnlib.Map(txts, func(tp gnml.TextPart) Chunk {
		res := Chunk{
			ItemID: pages[0].ItemID,
			Text:   tp.Content,
			Start:  uint(tp.Start),
			End:    uint(tp.End),
			Length: uint(tp.Length),
		}
		return res
	})
	chunks = findPages(pages, chunks)
	return chunks
}

func combinePages(pages []storage.Page) string {
	res := gnlib.Map(pages, func(p storage.Page) string {
		return p.Text
	})
	txt := strings.Join(res, "")
	return txt
}

func splitOverlap(itemID uint, txt string, size, overlap int) []Chunk {
	var chunks []Chunk

	for i := 0; i < len(txt); i += size - overlap {
		end := i + size
		if end > len(txt) {
			end = len(txt)
		}

		chunk := Chunk{
			ItemID: itemID,
			Text:   txt[i:end],
			// UUID:   gnuuid.New(txt).String(),
			Start: uint(i),
			End:   uint(end),
		}

		chunks = append(chunks, chunk)
	}
	return chunks
}

func findPages(pages []storage.Page, chunks []Chunk) []Chunk {
	res := make([]Chunk, len(chunks))
	pageIndex := 0

	for i, chunk := range chunks {
		var pagesForChunk []uint

		// Continue from the last page index
		for pageIndex < len(pages) {
			page := pages[pageIndex]

			// If the current page starts after the chunk ends, break the loop
			if page.Start > chunk.End {
				break
			}

			// Check if there is an overlap
			if page.End >= chunk.Start {
				// The chunk overlaps with the current page.
				pagesForChunk = append(pagesForChunk, page.ID)
			}

			// If the current page ends before the chunk ends, move to the next page
			if page.End < chunk.End {
				pageIndex++
				continue
			}

			// If the chunk ends within the current page, break the loop
			if chunk.End <= page.End {
				break
			}
		}
		chunk.PageIDs = pagesForChunk
		res[i] = chunk
	}

	return res
}
