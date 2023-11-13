package embedio

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/gnames/bhlquest/internal/io/dbshare"
	"github.com/gnames/bhlquest/pkg/config"
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
	var wg sync.WaitGroup
	wg.Add(1)

	go e.embedStream(chIn, chOut, &wg)
	go e.saveStream(chOut)

	e.loadChunks(chIn, itemIDs)
	wg.Wait()
	return nil
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

func (e *embedio) saveStream(chOut <-chan []text.Chunk) {
	for ch := range chOut {
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
