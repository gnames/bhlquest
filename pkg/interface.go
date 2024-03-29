package bhlquest

import (
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/output"
)

// BHLQuest provides functionality needed to apply AI to BHL.
type BHLQuest interface {
	// Init bootstraps AI engines providing necessary data and metadata.
	Init() error
	Ask(q string) (output.Answer, error)
	GetConfig() config.Config
	SetConfig(config.Config) BHLQuest
}
