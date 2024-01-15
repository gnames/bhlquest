package text

import "github.com/gnames/gnbhl/itembhl"

type Text interface {
	ItemToChunks(itemID uint) ([]Chunk, error)
	ItemByID(itemID uint) (itembhl.ItemBHL, error)
}
