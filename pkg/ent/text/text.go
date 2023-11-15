package text

import (
	"strings"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/storage"
	"github.com/gnames/gnlib"
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
	UUID      string
	Text      string
	Embedding []float32
	Distance  float64
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

func (t *text) pagesToChunks(pages []storage.Page) []Chunk {
	txt := combinePages(pages)
	chunks := splitOverlap(pages[0].ItemID, txt, 1000, 100)
	chunks = findPages(pages, chunks)
	return chunks
}

func combinePages(pages []storage.Page) string {
	res := gnlib.Map(pages, func(p storage.Page) string {
		return p.Text
	})
	txt := strings.Join(res, "\n")
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
		for pageIndex < len(pages) {
			page := pages[pageIndex]
			if chunk.End < page.Start {
				// The chunk ends before the current page starts; go to the next chunk.
				break
			}
			if chunk.Start <= page.End {
				// The chunk overlaps with the current page.
				pagesForChunk = append(pagesForChunk, page.ID)
				if chunk.End <= page.End {
					// The chunk ends within the current page; go to the next chunk.
					break
				}
			}
			pageIndex++ // Move to the next page.
		}
		chunk.PageIDs = pagesForChunk
		res[i] = chunk
	}

	return res
}
