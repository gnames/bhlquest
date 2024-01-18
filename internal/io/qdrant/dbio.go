package qdrant

import (
	"context"
	"log/slog"

	"github.com/gnames/bhlquest/internal/ent/text"
	pb "github.com/qdrant/go-client/qdrant"
)

var ctx = context.Background()

func (qd *qdrant) init() error {
	ctx := context.Background()
	_, err := qd.clientC.Delete(ctx, &pb.DeleteCollection{
		CollectionName: qd.cfg.QdrantDb,
	})
	if err != nil {
		return err
	}
	slog.Info(
		"Collection deleted",
		"collection",
		qd.cfg.QdrantDb)

	// Create new collection
	var defaultSegmentNumber uint64 = qd.cfg.QdrantSegmentsNum
	_, err = qd.clientC.Create(ctx, &pb.CreateCollection{
		CollectionName: qd.cfg.QdrantDb,
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
		qd.cfg.QdrantDb,
	)
	return nil
}

func (qd *qdrant) query(emb []float32) ([]text.Chunk, error) {
	ctx := context.Background()
	pts, err := qd.clientP.Search(ctx, &pb.SearchPoints{
		CollectionName: qd.cfg.QdrantDb,
		Vector:         emb,
		Limit:          25,
		// Include all payload and vectors in the search result
		WithVectors: &pb.WithVectorsSelector{
			SelectorOptions: &pb.WithVectorsSelector_Enable{Enable: true},
		},
		WithPayload: &pb.WithPayloadSelector{
			SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true},
		},
	})
	if err != nil {
		return nil, err
	}
	var res []text.Chunk
	for _, v := range pts.GetResult() {
		pl := v.Payload
		ch := text.Chunk{
			ID:        uint(v.Id.GetNum()),
			ItemID:    uint(pl["item_id"].GetIntegerValue()),
			Start:     uint(pl["chunk_offset"].GetIntegerValue()),
			Length:    uint(pl["chunk_length"].GetIntegerValue()),
			Embedding: v.Vectors.GetVector().GetData(),
			Score:     float64(v.GetScore()),
		}
		if ch.Score < qd.cfg.ScoreThreshold {
			continue
		}
		res = append(res, ch)
	}

	return res, nil
}

func (qd *qdrant) save(ch []text.Chunk) error {
	var res []*pb.PointStruct
	for _, v := range ch {
		res = append(res, getPoint(v))
	}
	return qd.upsertPoints(res)
}

func getPoint(ch text.Chunk) *pb.PointStruct {
	res := &pb.PointStruct{
		Id: &pb.PointId{PointIdOptions: &pb.PointId_Num{Num: uint64(ch.ID)}},
		Vectors: &pb.Vectors{
			VectorsOptions: &pb.Vectors_Vector{
				Vector: &pb.Vector{Data: ch.Embedding},
			},
		},

		Payload: map[string]*pb.Value{
			"item_id": {
				Kind: &pb.Value_IntegerValue{IntegerValue: int64(ch.ItemID)},
			},
			"chunk_offset": {
				Kind: &pb.Value_IntegerValue{IntegerValue: int64(ch.Start)},
			},
			"chunk_length": {
				Kind: &pb.Value_IntegerValue{IntegerValue: int64(ch.Length)},
			},
		},
	}
	return res

}

func (qd *qdrant) upsertPoints(pts []*pb.PointStruct) error {
	var err error
	// Upsert points
	waitUpsert := true
	_, err = qd.clientP.Upsert(
		ctx,
		&pb.UpsertPoints{
			CollectionName: qd.cfg.QdrantDb,
			Wait:           &waitUpsert,
			Points:         pts,
		},
	)
	return err
}
