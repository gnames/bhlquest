package embedio

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gnames/bhlquest/internal/io/dbshare"
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/embed"
	"github.com/jackc/pgx/v5/pgxpool"
)

type embedio struct {
	cfg config.Config
	db  *pgxpool.Pool
}

func New(cfg config.Config) (embed.Embed, error) {
	res := embedio{
		cfg: cfg,
	}
	db, err := pgxpool.New(
		context.Background(),
		dbshare.DbURL(cfg, cfg.DbBHLQuest),
	)
	if err != nil {
		return &res, err
	}
	res.db = db
	return &res, nil
}

func (e *embedio) Init() error {
	msg := fmt.Sprintf("Resetting '%s' database at '%s'.", e.cfg.DbBHLQuest, e.cfg.DbHost)
	slog.Info(msg)
	err := e.init()
	return err
}

func (e *embedio) Populate(itemIDs []int) error {
	return nil
}
