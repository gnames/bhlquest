package qdrant

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/answer"
	"github.com/gnames/bhlquest/pkg/ent/embed"
	"github.com/gnames/bhlquest/pkg/ent/llmutil"
	"github.com/gnames/bhlquest/pkg/ent/storage"
	"github.com/gnames/bhlquest/pkg/ent/text"
	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type qdrant struct {
	cfg      config.Config
	conn     *grpc.ClientConn
	clientC  pb.CollectionsClient
	clientP  pb.PointsClient
	txt      text.Text
	llm      llmutil.LlmUtil
	itemsNum int
}

func New(
	cfg config.Config,
	stg storage.Storage,
	llm llmutil.LlmUtil,
) (embed.Embed, error) {
	txt := text.New(cfg, stg)
	res := qdrant{
		cfg: cfg,
		txt: txt,
		llm: llm,
	}

	conn, err := grpc.Dial(cfg.QdrantHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	res.conn = conn
	res.clientC = pb.NewCollectionsClient(conn)
	res.clientP = pb.NewPointsClient(conn)

	return &res, nil
}

func (qd *qdrant) Init() error {
	msg := fmt.Sprintf(
		"Resetting '%s' database at '%s'.",
		qd.cfg.DbBHLQuest,
		qd.cfg.QdrantHost,
	)
	slog.Info(msg)
	err := qd.init()
	return err
}

func (qd *qdrant) LastItemID() uint {
	var res uint
	return res
}

func (qd *qdrant) SetItemsNum(itemsNum int) {
	qd.itemsNum = itemsNum
}

func (qd *qdrant) Populate(itemIDs []uint) error {
	chIn := make(chan []text.Chunk)
	chOut := make(chan []text.Chunk)
	var wg, saveWg sync.WaitGroup
	wg.Add(1)
	saveWg.Add(1)
	start := time.Now()

	go qd.embedStream(chIn, chOut, &wg)
	go qd.saveStream(chOut, start, &saveWg)

	qd.loadChunks(chIn, itemIDs)

	wg.Wait()
	close(chOut)

	saveWg.Wait()

	return nil
}

func (qd *qdrant) Embed(texts []string) ([][]float32, error) {
	return qd.llm.EmbedTexts(texts)
}

func (qd *qdrant) CrossEmbed(ss [][]string) ([]float64, error) {
	var res []float64
	return res, nil
}

func (qd *qdrant) Query(emb []float32) (answer.Answer, error) {
	var res answer.Answer
	chs, err := qd.query(emb)
	if err != nil {
		return res, err
	}
	var data []*answer.Result
	for i := range chs {
		var txt, txtExt string
		l := len(chs[i].PageIDs)

		txt, txtExt = qd.txt.ChunkText(chs[i])

		d := answer.Result{
			ChunkID:     chs[i].ID,
			ItemID:      chs[i].ItemID,
			PageIDStart: chs[i].PageIDs[0],
			PageIDEnd:   chs[i].PageIDs[l-1],
			Score:       chs[i].Score,
			Outlink: fmt.Sprintf(
				"%s%d",
				"https://www.biodiversitylibrary.org/page/",
				chs[i].PageIDs[0],
			),
			Text:    txt,
			TextExt: txtExt,
		}
		data = append(data, &d)
	}
	res.Results = data
	return res, nil
}

func (qd *qdrant) SetConfig(cfg config.Config) embed.Embed {
	copy := *qd
	copy.cfg = cfg
	return &copy
}
