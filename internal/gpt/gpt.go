package gpt

import (
	"fmt"
	"strings"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/output"
	"github.com/gnames/gnlib"
)

type gpt struct {
	cfg config.Config
	api API
}

func New(cfg config.Config, api API) GPT {
	res := gpt{
		cfg: cfg,
		api: api,
	}
	return &res
}

func (g *gpt) Summary(inp output.Answer) (string, error) {
	var res string
	if len(inp.Results) == 0 {
		return res, nil
	}

	texts := gnlib.Map(inp.Results, func(res *output.Result) string {
		return res.TextForSummary
	})
	if len(texts) >= 5 {
		texts = texts[:5]
	}
	data := strings.Join(texts, "\n\n")

	question := inp.Meta.Question
	userPrompt := fmt.Sprintf(Prompts["summary"], question, data)
	resp, err := g.api.Query(Prompts["system"], userPrompt)

	if err != nil {
		return res, err
	}

	res = resp.Choices[0].Message.Content
	return res, nil
}
