package rerank

import "github.com/gnames/bhlquest/pkg/ent/output"

type Reranker interface {
	// Rerank takes many potential answers and picks the best ones.
	Rerank(query string, texts []*output.Result) ([]*output.Result, error)
}
