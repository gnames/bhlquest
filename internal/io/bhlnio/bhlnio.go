package bhlnio

import (
	"context"

	"github.com/gnames/bhlquest/internal/io/dbshare"
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/bhln"
	"github.com/jackc/pgx/v5/pgxpool"
)

type bhlnio struct {
	cfg config.Config
	db  *pgxpool.Pool
}

func New(cfg config.Config) (bhln.BHLN, error) {
	res := bhlnio{
		cfg: cfg,
	}

	db, err := pgxpool.New(
		context.Background(),
		dbshare.DbURL(cfg, cfg.DbBHLNames),
	)
	if err != nil {
		return &res, err
	}
	res.db = db

	return &res, nil
}

func (bn *bhlnio) ItemIds(
	offset, limit int,
	taxa []string,
) ([]int, error) {
	return bn.dbItems(offset, limit, taxa)
}

func (bn *bhlnio) Close() {
	bn.db.Close()
}
