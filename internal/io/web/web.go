package web

import (
	"net/http"
	"strings"

	bhlquest "github.com/gnames/bhlquest/pkg"
	"github.com/gnames/bhlquest/pkg/ent/answer"
	"github.com/labstack/echo/v4"
)

type Duration struct {
	Total float32
}

// Data contains information needed to render web-pages.
type Data struct {
	Page           string
	Output         *answer.Answer
	Question       string
	Format         string
	FormatOptions  []string
	MaxResultsNum  int
	ScoreThreshold float64
	WithText       bool
	Version        string
}

type formInput struct {
	Question       string  `query:"question"`
	MaxResultsNum  int     `query:"max_results"`
	ScoreThreshold float64 `query:"score_threshold"`
	WithText       string  `query:"with_text"`
	Format         string  `query:"format"`
}

func homeGET(bq bhlquest.BHLQuest) func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{
			Page:           "home",
			Format:         "html",
			FormatOptions:  []string{"html", "json"},
			MaxResultsNum:  5,
			ScoreThreshold: 0.4,
			Version:        bhlquest.GetVersion().Version,
		}

		inp := new(formInput)
		err := c.Bind(inp)
		if err != nil {
			return err
		}

		if strings.TrimSpace(inp.Question) == "" {
			return c.Render(http.StatusOK, "layout", data)
		}

		cfg := bq.GetConfig()
		cfg.MaxResultsNum = inp.MaxResultsNum
		cfg.ScoreThreshold = inp.ScoreThreshold
		cfg.WithText = inp.WithText == "on"
		bq = bq.SetConfig(cfg)

		answ, err := bq.Ask(inp.Question)
		if err != nil {
			return err
		}
		data.Output = &answ
		data.Question = inp.Question
		data.Format = inp.Format
		data.MaxResultsNum = inp.MaxResultsNum
		data.ScoreThreshold = inp.ScoreThreshold
		data.WithText = inp.WithText == "on"

		switch data.Format {
		case "json":
			return c.JSON(http.StatusOK, data.Output)
		default:
			return c.Render(http.StatusOK, "layout", data)
		}
	}
}

func about(c echo.Context) error {
	data := Data{Page: "about", Version: bhlquest.GetVersion().Version}
	return c.Render(http.StatusOK, "layout", data)
}
