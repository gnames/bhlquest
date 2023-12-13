package gptio

import (
	"context"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/sashabaranov/go-openai"
)

type api struct {
	model  string
	client *openai.Client
}

func New(cfg config.Config) *api {
	res := api{
		model:  openai.GPT3Dot5Turbo1106,
		client: openai.NewClient(cfg.OpenaiAPIKey),
	}
	return &res
}

func (a *api) Query(system, prompt string) (openai.ChatCompletionResponse, error) {
	return a.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: a.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: system,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
}
