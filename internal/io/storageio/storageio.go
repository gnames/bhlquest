package storageio

import (
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/storage"
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
