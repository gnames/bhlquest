package bhlnio

import (
	"context"
	"fmt"

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
