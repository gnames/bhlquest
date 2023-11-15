package llmutilio

import (
	"net/http"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/llmutil"
	"github.com/gnames/bhlquest/pkg/ent/text"
	"github.com/gnames/gnfmt"
	"github.com/gnames/gnlib"
)

type llmutilio struct {
	cfg    config.Config
	url    string
	client *http.Client
	enc    gnfmt.Encoder
}

type embedPayload struct {
	Texts []string `json:"texts"`
}

func New(cfg config.Config) (llmutil.LlmUtil, error) {
	res := &llmutilio{
		cfg: cfg,
		enc: gnfmt.GNjson{},
	}
	err := res.conn()
	return res, err
}

func (l *llmutilio) Embed(
	cnks []text.Chunk,
) ([]text.Chunk, error) {
	res := make([]text.Chunk, len(cnks))
	texts := gnlib.Map(cnks, func(c text.Chunk) string {
		return c.Text
	})
	embs, err := l.embed(texts)
	if err != nil {
		return nil, err
	}
	for i := range embs {
		cnks[i].Embedding = embs[i]
		res[i] = cnks[i]
	}
	return res, nil
}

func (l *llmutilio) EmbedTexts(texts []string) ([][]float32, error) {
	return l.embed(texts)
}
