package llmutilio

import (
	"net/http"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/embed"
	"github.com/gnames/bhlquest/pkg/ent/llmutil"
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
	cnk []embed.Chunk,
) ([]embed.Chunk, error) {
	return nil, nil
}
