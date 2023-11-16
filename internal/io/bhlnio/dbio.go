package bhlnio

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (bn *bhlnio) dbItems(
	offset, limit int,
	taxa []string,
) ([]uint, error) {
	var rows pgx.Rows
	var err error
	var res []uint

	if len(bn.cfg.InitTaxa) > 0 {
		q := `
SELECT id
  FROM item_stats
  WHERE main_taxon = ANY($1::varchar[])`

		rows, err = bn.db.Query(context.Background(), q, bn.cfg.InitTaxa)
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
