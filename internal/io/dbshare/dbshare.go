package dbshare

import (
	"fmt"

	"github.com/gnames/bhlquest/pkg/config"
)

func DbURL(cfg config.Config, dbname string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DbUser, cfg.DbPass, cfg.DbHost, 5432, dbname)
}
