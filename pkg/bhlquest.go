package bhlquest

import (
	"fmt"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/bhln"
)

type bhlquest struct {
	cfg  config.Config
	bhln bhln.BHLN
}

func New(
	cfg config.Config,
	bhln bhln.BHLN,
) BHLQuest {
	res := bhlquest{
		cfg:  cfg,
		bhln: bhln,
	}
	return res
}

func (bq bhlquest) Init() error {
	ids, err := bq.bhln.ItemIds(0, 0, nil)
	if err != nil {
		return err
	}

	fmt.Println(ids)
	return nil
}

// GetVersion provides version information of the app.
func GetVersion() string {
	version := fmt.Sprintf("Version: %s\nBuild:   %s", Version, Build)
	return version
}
