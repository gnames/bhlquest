package qdrant

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gnames/bhlquest/internal/embed"
	"github.com/gnames/bhlquest/internal/ent/text"
	"github.com/gnames/bhlquest/internal/llmutil"
	"github.com/gnames/bhlquest/internal/storage"
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/output"
	"github.com/gnames/gnbhl/ent/pagebhl"
	"github.com/gnames/gnlib"
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

	conn, err := grpc.Dial(
		cfg.QdrantHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
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
		qd.cfg.QdrantDb,
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

// Query takes embedded question and resturns a slice of chunks that
// match putative answers.
func (qd *qdrant) Query(emb []float32) (output.Answer, error) {
	var res output.Answer
	chs, err := qd.query(emb)
	if err != nil {
		return res, err
	}
	var data []*output.Result
	for i := range chs {
		itm, err := qd.txt.ItemByID(chs[i].ItemID)
		if err != nil {
			return res, err
		}
		start := int(chs[i].Start)
		end := int(chs[i].Start + chs[i].Length)
		cnkBHL, err := itm.Chunk(start, end)
		if err != nil {
			return res, err
		}

		// get the text to sent to GPT for summary
		summaryText := cnkBHL.Text
		summaryChunk, err := itm.Chunk(start-500, start+2000)
		if err == nil {
			summaryText = summaryChunk.Text
		}

		var pagesText []string
		for _, v := range cnkBHL.Pages {
			pagesText = append(pagesText, markChunk(v, start, end))

		}

		d := output.Result{
			ChunkID: chs[i].ID,
			ItemID:  chs[i].ItemID,
			PageID:  cnkBHL.Pages[0].ID,
			Pages: gnlib.Map(itm.Pages(), func(p *pagebhl.PageBHL) output.Page {
				return output.Page{ID: p.ID, PageSeq: p.SeqNum}
			}),
			PageIndex: cnkBHL.PageIdxStart,

			Score: chs[i].Score,
			Outlink: fmt.Sprintf(
				"%s%d",
				"https://www.biodiversitylibrary.org/page/",
				cnkBHL.Pages[0].ID,
			),
			TextPages:      pagesText,
			TextForSummary: summaryText,
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

func markChunk(p *pagebhl.PageBHL, start, end int) string {
	offs := int(p.Offset)
	start = start - offs
	end = end - offs
	if (start < 0 && end < 0) || (start > len(p.Text) && end > len(p.Text)) {
		return p.Text
	}
	if start < 0 {
		start = 0
	}
	if end > len(p.Text) {
		end = len(p.Text)
	}
	var res string
	if start > 0 {
		res += p.Text[:start]
	}
	res += "<em>"
	res += p.Text[start:end]
	res += "</em>"
	if end < len(p.Text) {
		res += p.Text[end:]
	}
	return res
}
