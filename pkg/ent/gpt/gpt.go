package gpt

import (
	"fmt"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/answer"
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

func (g *gpt) Summary(inp answer.Answer) (string, error) {
	var res string
	if len(inp.Results) == 0 {
		return res, nil
	}
	data := inp.Results[0].PageText
	userPrompt := fmt.Sprintf(Prompts["summary"], data)
	resp, err := g.api.Query(Prompts["system"], userPrompt)

	if err != nil {
		return res, err
	}

	res = resp.Choices[0].Message.Content
	return res, nil
}
