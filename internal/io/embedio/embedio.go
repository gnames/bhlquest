package embedio

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gnames/bhlquest/internal/io/dbshare"
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/answer"
	"github.com/gnames/bhlquest/pkg/ent/embed"
	"github.com/gnames/bhlquest/pkg/ent/llmutil"
	"github.com/gnames/bhlquest/pkg/ent/storage"
	"github.com/gnames/bhlquest/pkg/ent/text"
	"github.com/gnames/gnfmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type embedio struct {
	cfg      config.Config
	db       *pgxpool.Pool
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
	res := embedio{
		cfg: cfg,
		txt: txt,
		llm: llm,
	}

	db, err := pgxpool.New(
		context.Background(),
		dbshare.DbURL(cfg, cfg.DbBHLQuest),
	)
	if err != nil {
		return &res, err
	}
	res.db = db
	return &res, nil
}

func (e *embedio) Init() error {
	msg := fmt.Sprintf("Resetting '%s' database at '%s'.", e.cfg.DbBHLQuest, e.cfg.DbHost)
	slog.Info(msg)
	err := e.init()
	return err
}

func (e *embedio) SetItemsNum(i int) {
	e.itemsNum = i
}

func (e *embedio) Populate(itemIDs []uint) error {
	chIn := make(chan []text.Chunk)
	chOut := make(chan []text.Chunk)
	var wg, saveWg sync.WaitGroup
	wg.Add(1)
	saveWg.Add(1)
	start := time.Now()

	go e.embedStream(chIn, chOut, &wg)
	go e.saveStream(chOut, start, &saveWg)

	e.loadChunks(chIn, itemIDs)

	wg.Wait()
	close(chOut)

	saveWg.Wait()

	slog.Info("Creating index for embedding field.")
	return e.createIndex()
}

func (e *embedio) LastItemID() uint {
	return e.lastItemID()
}

func (e *embedio) Embed(texts []string) ([][]float32, error) {
	return e.llm.EmbedTexts(texts)
}

func (e *embedio) Query(emb []float32) (answer.Answer, error) {
	var res answer.Answer
	chs, err := e.query(emb)
	if err != nil {
		return res, err
	}
	var data []answer.Result
	for i := range chs {
		var txt, txtExt string
		l := len(chs[i].PageIDs)

		if e.cfg.WithText {
			txt, txtExt = e.txt.ChunkText(chs[i])
		}

		d := answer.Result{
			ItemID:      chs[i].ItemID,
			PageIDStart: chs[i].PageIDs[0],
			PageIDEnd:   chs[i].PageIDs[l-1],
			Score:       1 - chs[i].Distance,
			Outlink: fmt.Sprintf(
				"%s%d",
				"https://www.biodiversitylibrary.org/page/",
				chs[i].PageIDs[0],
			),
			Text:    txt,
			TextExt: txtExt,
		}
		data = append(data, d)
	}
	res.Results = data
	return res, nil
}

func (e *embedio) SetConfig(cfg config.Config) embed.Embed {
	copy := *e
	copy.cfg = cfg
	return &copy
}

func (e *embedio) embedStream(
	chIn <-chan []text.Chunk,
	chOut chan<- []text.Chunk,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for chnks := range chIn {
		chnks, err := e.llm.Embed(chnks)
		if err != nil {
			slog.Error("Cannot embed with llmutil", "error", err)
			os.Exit(1)
		}
		chOut <- chnks
	}
}

func (e *embedio) saveStream(
	chOut <-chan []text.Chunk,
	start time.Time,
	saveWg *sync.WaitGroup,
) {
	defer saveWg.Done()
	var count int
	for ch := range chOut {
		count = incrLog(start, e.itemsNum, count, int(ch[0].ItemID), 1)
		e.save(ch)

	}
	fmt.Fprint(os.Stderr, "\r")
}

func incrLog(start time.Time, total, count, itemID, incr int) int {
	count += incr
	if count%500 == 0 {
		fmt.Fprint(os.Stderr, "\r")
		slog.Info(logStr(start, total, count), "ItemID", itemID)
	} else if count%10 == 0 {
		fmt.Fprintf(os.Stderr, "\r%s", strings.Repeat(" ", 80))
		fmt.Fprint(os.Stderr, "\r"+logStr(start, total, count))
	}
	return count
}

func logStr(start time.Time, total, count int) string {
	rate := itemsPerHour(start, count)
	countStr := humanize.Comma(int64(count))
	perHourStr := humanize.Comma(int64(rate))
	percent := 100 * float64(count) / float64(total)
	eta := 3600 * float64(total-count) / rate
	etaStr := gnfmt.TimeString(eta)
	return fmt.Sprintf(
		"%s items (%0.1f%%), %s items/hr: ETA %s",
		countStr, percent, perHourStr, etaStr,
	)
}

func itemsPerHour(start time.Time, count int) float64 {
	dur := float64(time.Since(start)) / float64(time.Hour)
	return float64(count) / dur
}

func (e *embedio) loadChunks(
	chIn chan<- []text.Chunk,
	itemIDs []uint,
) {
	for _, id := range itemIDs {
		chunks, err := e.txt.TextToChunks(id)
		if err != nil {
			slog.Error("Cannot create chunks", "error", err)
			os.Exit(1)
		}
		chIn <- chunks
	}
	close(chIn)
}
