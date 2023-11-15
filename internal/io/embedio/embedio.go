package embedio

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"

	"github.com/gnames/bhlquest/internal/io/dbshare"
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/bhlquest/pkg/ent/answer"
	"github.com/gnames/bhlquest/pkg/ent/embed"
	"github.com/gnames/bhlquest/pkg/ent/llmutil"
	"github.com/gnames/bhlquest/pkg/ent/storage"
	"github.com/gnames/bhlquest/pkg/ent/text"
	"github.com/jackc/pgx/v5/pgxpool"
)

type embedio struct {
	cfg config.Config
	db  *pgxpool.Pool
	txt text.Text
	llm llmutil.LlmUtil
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

func (e *embedio) Populate(itemIDs []uint) error {
	chIn := make(chan []text.Chunk)
	chOut := make(chan []text.Chunk)
	var wg, saveWg sync.WaitGroup
	wg.Add(1)
	saveWg.Add(1)

	go e.embedStream(chIn, chOut, &wg)
	go e.saveStream(chOut, &saveWg)

	e.loadChunks(chIn, itemIDs)

	wg.Wait()
	close(chOut)

	saveWg.Wait()

	return nil
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
		l := len(chs[i].PageIDs)
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
		}
		data = append(data, d)
	}
	res.Results = data
	return res, nil
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
			slog.Error(err.Error())
			os.Exit(1)
		}
		chOut <- chnks
	}
}

func (e *embedio) saveStream(chOut <-chan []text.Chunk, saveWg *sync.WaitGroup) {
	defer saveWg.Done()
	var count int
	for ch := range chOut {
		count++
		slog.Info(strconv.Itoa(count))
		e.save(ch)
		_ = ch
	}
}

func (e *embedio) loadChunks(
	chIn chan<- []text.Chunk,
	itemIDs []uint,
) {
	for _, id := range itemIDs {
		chunks, err := e.txt.TextToChunks(id)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		chIn <- chunks
	}
	close(chIn)
}
