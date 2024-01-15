package qdrant

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gnames/bhlquest/internal/ent/text"
	"github.com/gnames/gnfmt"
)

func (qd qdrant) embedStream(
	chIn <-chan []text.Chunk,
	chOut chan<- []text.Chunk,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for chnks := range chIn {
		chnks, err := qd.llm.Embed(chnks)
		if err != nil {
			slog.Error("Cannot embed with llmutil", "error", err)
			os.Exit(1)
		}
		chOut <- chnks
	}
}

func (qd qdrant) saveStream(
	chOut <-chan []text.Chunk,
	start time.Time,
	saveWg *sync.WaitGroup,
) {
	defer saveWg.Done()
	var count int
	for ch := range chOut {
		count = incrLog(start, qd.itemsNum, count, int(ch[0].ItemID), 1)
		qd.save(ch)
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

func (qd *qdrant) loadChunks(
	chIn chan<- []text.Chunk,
	itemIDs []uint,
) {
	var count uint
	for _, id := range itemIDs {
		chunks, err := qd.txt.ItemToChunks(id)
		if err != nil {
			slog.Error("Cannot create chunks", "error", err)
			os.Exit(1)
		}
		for i := range chunks {
			count++
			chunks[i].ID = count
		}
		chIn <- chunks
	}
	close(chIn)
}
