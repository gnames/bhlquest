package text_test

import (
	"testing"

	"github.com/gnames/bhlquest/internal/ent/text"
	"github.com/gnames/bhlquest/internal/io/storageio"
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/stretchr/testify/assert"
)

func getText(t *testing.T) text.Text {
	cfg := config.New(config.OptBHLDir("../../testdata/bhl"))
	stg := storageio.New(cfg)
	txt := text.New(cfg, stg)
	return txt
}

func TestItemToChunks(t *testing.T) {
	assert := assert.New(t)
	txt := getText(t)
	chunks, err := txt.ItemToChunks(100100)
	assert.Nil(err)
	for _, c := range chunks {
		assert.Greater(len(c.Text), 100)
	}
	assert.Greater(len(chunks), 10)
}

func TestItemByID(t *testing.T) {
	assert := assert.New(t)
	txt := getText(t)

	itm, err := txt.ItemByID(uint(100100))
	assert.Nil(err)
	assert.Equal(68, len(itm.Pages()))
	assert.Contains(itm.Text(), "NATIONAL FISH HATCHERY SYSTEM")
	cnk, err := itm.Chunk(1000, 2000)
	assert.Nil(err)
	assert.Equal(4, len(cnk.Pages))
	assert.Equal(844, len(cnk.Pages[0].Text))
	assert.Contains(cnk.Text, "NATIONAL FISH HATCHERY SYSTEM")
}
