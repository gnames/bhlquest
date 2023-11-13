package bhlnio

import "context"

func (bn *bhlnio) dbItems(
	offset, limit int,
	taxa []string,
) ([]uint, error) {
	q := `
SELECT id
  FROM item_stats
  WHERE main_taxon in ('Aves')
`
	var res []uint
	rows, err := bn.db.Query(context.Background(), q)

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
