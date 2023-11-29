package embedio

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gnames/bhlquest/pkg/ent/text"
	pgvector "github.com/pgvector/pgvector-go"
)

func (e *embedio) init() error {
	cxt := context.Background()
	q := `
DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO %s;
CREATE EXTENSION IF NOT EXISTS vector;
COMMENT ON SCHEMA public IS 'standard public schema'`
	q = fmt.Sprintf(q, e.cfg.DbUser)
	_, err := e.db.Exec(cxt, q)
	if err != nil {
		return fmt.Errorf("-> db.Exec %w", err)
	}
	slog.Info("Creating tables.")
	q = `
CREATE TABLE chunks (
	id bigserial PRIMARY KEY,
	item_id int,
	page_id bigint,
	page_id_end bigint,
	item_offset int,
	embedding vector(768)
	)		
	`
	_, err = e.db.Exec(cxt, q)
	if err != nil {
		err = fmt.Errorf("-> migrate %w", err)
	}
	return err

}

func (e *embedio) lastItemID() uint {
	ctx := context.Background()
	q := `SELECT item_id FROM chunks WHERE id = (SELECT MAX(id) FROM chunks)`
	row := e.db.QueryRow(ctx, q)
	var id uint
	err := row.Scan(&id)
	if err != nil {
		return 0
	}
	return id
}

func (e *embedio) save(chunks []text.Chunk) error {
	ctx := context.Background()
	q := `INSERT INTO chunks
		(item_id, page_id, page_id_end, item_offset, embedding)
		VALUES ($1, $2, $3, $4, $5)`
	tx, err := e.db.Begin(ctx)
	if err != nil {
		return err
	}
	for _, v := range chunks {
		pIDs := v.PageIDs
		l := len(pIDs)
		if l == 0 {
			continue
		}
		vec := pgvector.NewVector(v.Embedding)
		_, err := e.db.Exec(ctx, q, v.ItemID, pIDs[0], pIDs[l-1], v.Start, vec)
		if err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (e *embedio) createIndex() error {
	ctx := context.Background()
	q := `
CREATE INDEX
	ON chunks
	USING ivfflat (embedding vector_cosine_ops)
	WITH (lists = 3000);
`

	_, err := e.db.Exec(ctx, q)
	return err
}

func (e *embedio) query(emb []float32) ([]text.Chunk, error) {
	ctx := context.Background()
	q := `
SELECT id, item_id, page_id, page_id_end, item_offset,
       embedding <=> $1 AS dot
  FROM chunks
  WHERE (embedding <=> $1) < $2
  ORDER BY dot
  LIMIT $3
`
	rows, err := e.db.Query(ctx, q,
		pgvector.NewVector(emb),
		e.cfg.ScoreThreshold,
		e.cfg.MaxResultsNum,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var res []text.Chunk
	for rows.Next() {
		var ch text.Chunk
		var start, end int
		var dotProd float64
		err = rows.Scan(&ch.ID, &ch.ItemID, &start, &end, &ch.Start, &dotProd)
		if err != nil {
			panic(err)
		}
		ch.PageIDs = append(ch.PageIDs, uint(start))
		if start != end {
			ch.PageIDs = append(ch.PageIDs, uint(end))
		}
		ch.Distance = dotProd
		res = append(res, ch)
	}

	if rows.Err() != nil {
		panic(rows.Err())
	}

	return res, nil
}
