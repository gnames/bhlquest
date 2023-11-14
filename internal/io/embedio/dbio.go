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
	embedding vector(384)
	)		
	`
	_, err = e.db.Exec(cxt, q)
	if err != nil {
		err = fmt.Errorf("-> migrate %w", err)
	}
	return err

}

func (e *embedio) save(chunks []text.Chunk) error {
	cxt := context.Background()
	q := `INSERT INTO chunks
		(item_id, page_id, page_id_end, embedding)
		VALUES ($1, $2, $3, $4)`
	for _, v := range chunks {
		pIDs := v.PageIDs
		l := len(pIDs)
		if l == 0 {
			continue
		}
		vec := pgvector.NewVector(v.Embedding)
		_, err := e.db.Exec(cxt, q, v.ItemID, pIDs[0], pIDs[l-1], vec)
		if err != nil {
			return err
		}
	}
	return nil
}
