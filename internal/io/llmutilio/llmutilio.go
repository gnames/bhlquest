package llmutilio

import (
	"net/http"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/llmutil"
	"github.com/gnames/bhlquest/pkg/ent/text"
)

type llmutilio struct {
	cfg    config.Config
	url    string
	client *http.Client
}

func New(cfg config.Config) (llmutil.LlmUtil, error) {
	res := &llmutilio{
		cfg: cfg,
	}
	err := res.conn()
	return res, err
}

func (l *llmutilio) Embed(
	cnk []text.Chunk,
) ([]text.Chunk, error) {
	return nil, nil
}
