package qdrant

import (
	"context"
	"log/slog"

	"github.com/gnames/bhlquest/pkg/ent/text"
	pb "github.com/qdrant/go-client/qdrant"
)

func (qd *qdrant) init() error {
	ctx := context.Background()
	_, err := qd.clientC.Delete(ctx, &pb.DeleteCollection{
		CollectionName: qd.cfg.DbBHLQuest,
	})
	if err != nil {
		return err
	}
	slog.Info(
		"Collection deleted",
		"collection",
		qd.cfg.DbBHLQuest)

	// Create new collection
	var defaultSegmentNumber uint64 = qd.cfg.QdrantSegmentsNum
	_, err = qd.clientC.Create(ctx, &pb.CreateCollection{
		CollectionName: qd.cfg.DbBHLQuest,
		VectorsConfig: &pb.VectorsConfig{Config: &pb.VectorsConfig_Params{
			Params: &pb.VectorParams{
				Size:     qd.cfg.VectorSize,
				Distance: pb.Distance_Cosine,
			},
		}},
		OptimizersConfig: &pb.OptimizersConfigDiff{
			DefaultSegmentNumber: &defaultSegmentNumber,
		},
	})
	if err != nil {
		return err
	}
	slog.Info(
		"Collection created",
		"collection",
		qd.cfg.DbBHLQuest,
	)
	return nil
}

func (qd *qdrant) query(emb []float32) ([]text.Chunk, error) {
	ctx := context.Background()
	pts, err := qd.clientP.Search(ctx, &pb.SearchPoints{
		CollectionName: qd.cfg.DbBHLQuest,
		Vector:         emb,
		Limit:          25,
		// Include all payload and vectors in the search result
		WithVectors: &pb.WithVectorsSelector{SelectorOptions: &pb.WithVectorsSelector_Enable{Enable: true}},
		WithPayload: &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true}},
	})
	if err != nil {
		return nil, err
	}
	var res []text.Chunk
	for _, v := range pts.GetResult() {
		// fmt.Printf("POINT: %#v\n", v)
		pl := v.Payload
		pgSt := uint(pl["page_id"].GetIntegerValue())
		pgEnd := uint(pl["page_end_id"].GetIntegerValue())
		ch := text.Chunk{
			ID:        uint(v.Id.GetNum()),
			ItemID:    uint(pl["item_id"].GetIntegerValue()),
			PageIDs:   []uint{pgSt, pgEnd},
			Start:     uint(pl["item_offset"].GetIntegerValue()),
			Embedding: v.Vectors.GetVector().GetData(),
		}
		res = append(res, ch)
	}

	return res, nil
}
