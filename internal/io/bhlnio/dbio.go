package bhlnio

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gnames/bhlquest/internal/ent/ref"
	"github.com/jackc/pgx/v5"
)

func (bn *bhlnio) dbItems(
	offset, limit int,
	taxa []string,
) ([]uint, error) {
	var rows pgx.Rows
	var err error
	var res []uint

	table := "items"
	if len(bn.cfg.InitClasses) > 0 || len(bn.cfg.InitTaxa) > 0 {
		table = "item_stats"
	}

	q := fmt.Sprintf("SELECT id FROM %s", table)

	if table == "item_stats" {
		q += `
WHERE main_class = ANY($1::varchar[]) OR
      main_taxon = ANY($2::varchar[])
`

		rows, err = bn.db.Query(
			context.Background(),
			q,
			bn.cfg.InitClasses,
			bn.cfg.InitTaxa,
		)
	} else {
		q := `SELECT id FROM items`
		rows, err = bn.db.Query(context.Background(), q)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id uint
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

func (bn *bhlnio) dbPages(itemID uint) (map[uint]string, error) {
	q := `SELECT id, page_num FROM pages WHERE item_id = $1`
	rows, err := bn.db.Query(context.Background(), q, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make(map[uint]string)
	for rows.Next() {
		var id uint
		var pageNum sql.NullInt32
		err = rows.Scan(&id, &pageNum)
		if err != nil {
			return nil, err
		}
		if pageNum.Valid {
			res[id] = strconv.Itoa(int(pageNum.Int32))
		} else {
			res[id] = ""
		}
	}
	return res, nil
}

func (bn *bhlnio) dbReference(pageIDs []int) (map[int]ref.Reference, error) {
	res := make(map[int]ref.Reference)
	q := `
SELECT p.id, i.id, i.title_name, i.vol, i.title_doi, i.title_year_start, 
  i.title_year_end, i.title_lang,
  p.page_num
  FROM pages p 
    JOIN items i 
      ON i.id = p.item_id
    WHERE p.id = any($1::int[])
`
	rows, err := bn.db.Query(context.Background(), q, pageIDs)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var r ref.Reference
		err = rows.Scan(&r.PageID, &r.ItemID, &r.TitleName, &r.Volume,
			&r.TitleDOI, &r.TitleYearStart, &r.TitleYearEnd, &r.TitleLang,
			&r.PageNumber,
		)
		if err != nil {
			return res, err
		}
		r.Fingerprint = r.GetFingerprint()
		res[r.PageID] = r
	}
	return res, nil
}
