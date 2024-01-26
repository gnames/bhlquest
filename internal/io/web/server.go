package web

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
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
// @description This API serves the BHLQuest app. It locates relevant sections in the Biodiversity Heritage Library that correspond to a user's query. \n\nCode repository: https://github.com/gnames/bhlquest. \n\nAccess the API on the production server: https://bhlquest.globalnames.org/api/v1.

// @contact.name Dmitry Mozzherin
// @contact.url https://github.com/dimus
// @contact.email dmozzherin@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// Server Definitions
// @Server http://localhost:8555 Description for local server
// @Server https://bhlquest.globalnames.org Description for production server

// @host localhost:8555
// @host bhlquest.globalnames.org
// @BasePath /api/v1
func Run(bq bhlquest.BHLQuest) {
	var err error
	port := bq.GetConfig().Port
	slog.Info("Starting HTTP API server", "port", port)
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	setLogger(e)

	e.Renderer, err = NewTemplate()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/", homeGET(bq))
	e.GET("/about", about)
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

	fs := http.FileServer(http.FS(static))
	e.GET("/static/*", echo.WrapHandler(fs))

	e.Logger.Fatal(e.StartServer(s))
}

func setLogger(e *echo.Echo) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
}
