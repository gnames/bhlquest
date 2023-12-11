package gpt

import (
	"github.com/gnames/bhlquest/pkg/ent/answer"
	"github.com/sashabaranov/go-openai"
)

type GPT interface {
	Summary(answer.Answer) (string, error)
}

type API interface {
	Query(system, prompt string) (openai.ChatCompletionResponse, error)
}
