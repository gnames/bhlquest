package embedio

import (
	"context"
	"fmt"
	"log/slog"
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
	title_id int,
	page_id int,
	uuid uuid,
	embedding vector(384)
	)		
	`
	_, err = e.db.Exec(cxt, q)
	if err != nil {
		err = fmt.Errorf("-> migrate %w", err)
	}
	return err

}
