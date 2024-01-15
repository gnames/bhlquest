package llmutil

import "github.com/gnames/bhlquest/internal/ent/text"

type LlmUtil interface {
	Embed([]text.Chunk) ([]text.Chunk, error)
	EmbedTexts([]string) ([][]float32, error)
	CrossEmbedPairs([][]string) ([]float64, error)
}
