package web

import (
	"net/http"
	"net/url"
	"strconv"

	bhlquest "github.com/gnames/bhlquest/pkg" // for docs
	"github.com/labstack/echo/v4"
)

// info gives information where to find docs.
// @Summary Information about the API documentation
// @Description Gives information where to find docs.
// @ID get-info
// @Produce plain
// @Success 200 {string} string "API documentation URL"
// @Router / [get]
func info(c echo.Context) error {
	return c.String(http.StatusOK,
		`The API is described at

https://bhlquest.globalnames.org/apidoc/
`,
	)
}

// ping checks if the API is online
// @Summary Check API status
// @Description Checks if the API is online and returns a simple response if it is.
// @ID get-ping
// @Produce plain
// @Success 200 {string} string "API status response"
// @Router /ping [get]
func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

// version returns back the version of BHLQuest
// @Summary Get BHLQuest version
// @Description Retrieves the current version of the BHLQuest application.
// @ID get-version
// @Produce json
// @Success 200 {object} gnvers.Version "Successful response with version information"
// @Router /version [get]
func ver(c echo.Context) error {
	result := bhlquest.GetVersion()
	return c.JSON(http.StatusOK, result)
}

// ask receives a question and returns back a list of pages containing the
// answer.
// @Summary Ask a question
// @Description This endpoint receives a question about BHL and returns a list of pages containing the answer.
// @ID ask-question
// @Produce json
// @Param question path string true "A question to ask BHL about."
// @Param max-results query integer false "The maximum number or returned results."
// @Param score-threshold query number false "A score threshold from 0.0 to 1.0"
// @Param with-text query bool false "Shows matched text in results"
// @Success 200 {array} answer.Answer "List of pages containing the answer"
// @Router /ask/{question} [get]
func ask(bq bhlquest.BHLQuest) func(c echo.Context) error {
	return func(c echo.Context) error {
		q, _ := url.QueryUnescape(c.Param("question"))
		cfg := bq.GetConfig()

		maxStr, _ := url.QueryUnescape(c.QueryParam("max-results"))
		max, err := strconv.Atoi(maxStr)
		if err == nil && max > 0 && max < 100 {
			cfg.MaxResultsNum = max
		}

		thrStr, _ := url.QueryUnescape(c.QueryParam("score-threshold"))
		thr, err := strconv.ParseFloat(thrStr, 64)
		if err == nil && thr >= 0.0 && thr <= 1.0 {
			cfg.ScoreThreshold = 1 - thr
		}

		bq = bq.SetConfig(cfg)

		answ, err := bq.Ask(q)
		cfg = bq.GetConfig()
		answ.MaxResultsNum = cfg.MaxResultsNum
		answ.ScoreThreshold = 1 - cfg.ScoreThreshold
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, answ)
	}
}
