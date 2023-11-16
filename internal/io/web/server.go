package web

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/gnames/bhlquest/docs"
	bhlquest "github.com/gnames/bhlquest/pkg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var (
	apiPath = "/api/v1/"
)

//go:embed static
var static embed.FS

// @title BHLQuest API
// @version 1.0
// @description This is API to BHLQuest app. It finds places in Biodiversity Heritage Library that correspond to an asked question. \n\nCode: https://github.com/gnames/bhlquest. \n\nProduction server: https://bhlquest.globalnames.org/api/v1

// @contact.name Dmitry Mozzherin
// @contact.url https://github.com/dimus
// @contact.email dmozzherin@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// Server Definitions
// @Server http://localhost:8555 Description for local server
// @Server https://bhlquest.globalnames.org Description for production server

// @host localhost:8555
// @BasePath /api/v1
func Run(bq bhlquest.BHLQuest) {
	port := bq.GetConfig().Port
	slog.Info("Starting HTTP API server", "port", port)
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	e.GET("/", home(bq))
	e.GET("/apidoc/*", echoSwagger.WrapHandler)
	e.GET("/api", info)
	e.GET("/api/", info)
	e.GET("/api/v1", info)
	e.GET(apiPath, info)
	e.GET(apiPath, info)
	e.GET(apiPath+"ping", ping)
	e.GET(apiPath+"version", ver)
	e.GET(apiPath+"ask/:question", ask(bq))

	addr := fmt.Sprintf(":%d", port)
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}
