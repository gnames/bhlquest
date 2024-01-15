package storageio

import (
	"github.com/gnames/bhlquest/internal/storage"
	"github.com/gnames/bhlquest/pkg/config"
)

type storageio struct {
	cfg config.Config
}

func New(cfg config.Config) storage.Storage {
	res := &storageio{
		cfg: cfg,
	}
	return res
}
