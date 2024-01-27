package bhln

import "github.com/gnames/bhlquest/internal/ent/ref"

// BHLN takes data from BHLnames.
type BHLN interface {
	// References returns back bibliographic references for a given
	// list of page ids.
	References(pageIDs []int) (map[int]ref.Reference, error)

	// PageNums returns back a map of page ids to page numbers for a given
	// BHL item.
	PageNums(itemID uint) (map[uint]uint, error)

	// ItemIds retuns back ids of BHL items that need to be embedded.
	ItemIds(offset, limit int, taxa []string) ([]uint, error)

	// Close cleans up database connections
	Close()
}
