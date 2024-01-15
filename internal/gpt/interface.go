package gpt

import (
	"github.com/gnames/bhlquest/pkg/ent/output"
	openai "github.com/sashabaranov/go-openai"
)

// GPT is an interface for GPT-3.5 API. It provides
type GPT interface {
	Summary(output.Answer) (string, error)
}

// API is
type API interface {
	Query(system, prompt string) (openai.ChatCompletionResponse, error)
}
