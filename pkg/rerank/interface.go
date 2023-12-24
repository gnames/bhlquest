package rerank

import (
	"github.com/gnames/bhlquest/pkg/ent/answer"
)

type Reranker interface {
	// Rerank takes many potential answers and picks the best ones.
	Rerank(query string, texts []answer.Result) ([]answer.Result, error)
}
