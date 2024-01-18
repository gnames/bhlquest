package ref

import (
	"database/sql"
	"fmt"

	"github.com/gnames/bhlnames/pkg/ent/abbr"
)

// Reference is a struct that represents a reference to a page in BHL.
type Reference struct {
	// ItemID is the ID of a BHL Item, such as a book.
	ItemID int

	// PageID is the ID of a BHL Page.
	PageID int

	// Fingerprint helps to identify duplicates in results.
	// It is generated by simplified version of the title, volume,
	// and page_number
	Fingerprint string

	// TitleName is the name of the book or journal.
	// A BHL item can be a part of a Title, for example a journal
	// volume/volumes are part of a journal.
	TitleName string

	// TitleDOI is the DOI of the book or journal.
	TitleDOI string

	// TitleLang is the language of the book or journal.
	TitleLang string

	// TitleYearStart is the year of publication of the book or journal.
	// In case of a journal, it's the year of the first volume.
	TitleYearStart sql.NullInt16

	// TitleYearEnd is the year of the last volume of a journal.
	TitleYearEnd sql.NullInt16

	// Volume is the number of the volume if an item is a volume of a journal.
	Volume string

	// PageNumber is the number of the page that appear in the print
	// of the publication.
	PageNumber sql.NullInt16
}

// String returns a string representation of the reference.
func (r Reference) String() string {
	res := r.TitleName
	if r.TitleYearStart.Valid {
		res += r.yearString()
	}
	if r.Volume != "" {
		res += ", " + r.Volume
	}
	if r.PageNumber.Valid {
		res += ", p." + fmt.Sprint(r.PageNumber.Int16)
	}
	return res
}

// yearString returns a string representation of the year of publication.
func (r Reference) yearString() string {
	if r.TitleYearEnd.Valid {
		return fmt.Sprintf(
			" (%d-%d)",
			r.TitleYearStart.Int16,
			r.TitleYearEnd.Int16,
		)
	}
	return fmt.Sprintf(" (%d)", r.TitleYearStart.Int16)
}

func (r Reference) GetFingerprint() string {
	titleAbbr := abbr.AbbrMax(r.TitleName, nil)
	page := int(r.PageNumber.Int16)
	return fmt.Sprintf("%s|%d", titleAbbr, page)
}