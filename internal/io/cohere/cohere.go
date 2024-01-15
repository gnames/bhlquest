package cohere

import (
	"cmp"
	"context"
	"slices"

	"github.com/gnames/bhlquest/internal/rerank"
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/output"
	"github.com/gnames/gnlib"
	coherego "github.com/henomis/cohere-go"
	"github.com/henomis/cohere-go/request"
	"github.com/henomis/cohere-go/response"
)

type cohere struct {
	key    string
	client *coherego.Client
}

func New(cfg config.Config) rerank.Reranker {
	res := cohere{
		key:    cfg.CohereAPIKey,
		client: coherego.New(cfg.CohereAPIKey),
	}
	return &res
}

func (c *cohere) Rerank(
	query string,
	rs []*output.Result,
) ([]*output.Result, error) {
	maxChunksPerDoc := 10

	resp := &response.Rerank{}
	txts := gnlib.Map(rs, func(r *output.Result) string {
		return r.TextForSummary
	})
	req := request.Rerank{
		ReturnDocuments: true,
		MaxChunksPerDoc: &maxChunksPerDoc,
		Query:           query,
		Documents:       txts,
	}

	err := c.client.Rerank(context.Background(), &req, resp)
	if err != nil {
		return nil, err
	}

	for i := range resp.Results {
		rs[i].CrossScore = resp.Results[i].RelevanceScore
	}
	slices.SortFunc(rs, func(a, b *output.Result) int {
		return cmp.Compare(b.CrossScore, a.CrossScore)
	})
	return rs, nil
}
