package web

import (
	"net/http"

	bhlquest "github.com/gnames/bhlquest/pkg"
	"github.com/gnames/bhlquest/pkg/ent/answer"
	"github.com/labstack/echo/v4"
)

type Duration struct {
	Total float32
}

// Data contains information needed to render web-pages.
type Data struct {
	Input       string
	Output      answer.Answer
	Duration    Duration
	Page        string
	Format      string
	UniqueNames bool
	Version     string
	APIDoc      string
}

func home(bq bhlquest.BHLQuest) func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{
			Page:    "home",
			Version: bhlquest.GetVersion().Version}
		return c.Render(http.StatusOK, "layout", data)
	}
}

func apidoc(bq bhlquest.BHLQuest) func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{
			Page:    "apidoc",
			Version: bhlquest.GetVersion().Version,
			APIDoc:  bq.GetConfig().APIDocURL,
		}
		return c.Render(http.StatusOK, "layout", data)
	}
}
