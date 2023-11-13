package llmutil

import "github.com/gnames/bhlquest/pkg/ent/embed"

type LlmUtil interface {
	Embed([]embed.Chunk) []embed.Chunk
}
