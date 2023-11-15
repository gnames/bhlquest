package llmutil

import "github.com/gnames/bhlquest/pkg/ent/text"

type LlmUtil interface {
	Embed([]text.Chunk) ([]text.Chunk, error)
	EmbedTexts([]string) ([][]float32, error)
}
