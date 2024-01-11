package bhlnio

import (
	"context"
	"fmt"

	"github.com/gnames/bhlquest/pkg/ent/ref"
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
		q := `select id from items`
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
