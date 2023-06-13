package cmd

import (
	"fmt"
	"net/url"
	"time"

	"github.com/jtarchie/forum/cache"
	"github.com/jtarchie/forum/db"
	"github.com/jtarchie/forum/services"
	"github.com/jtarchie/forum/templates"
	"github.com/jtarchie/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ServerCmd struct {
	Port     int      `help:"port to run http server"`
	DBServer *url.URL `help:"URL to the rqlite API" default:"http://localhost:4001"`
}

func (c *ServerCmd) Run() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("could not create logger: %w", err)
	}

	e := echo.New()
	e.Use(middleware.ZapLogger(logger))

	// run migrations
	client, err := db.NewClient(c.DBServer.String())
	if err != nil {
		return fmt.Errorf("could create client: %w", err)
	}

	err = services.Migration(client, logger)
	if err != nil {
		return fmt.Errorf("could not migrate: %w", err)
	}

	cachedForums := cache.NewFunc(services.ListForums, time.Minute)

	e.GET("/", func(c echo.Context) error {
		forums, err := cachedForums.Invoke(client)
		if err != nil {
			return fmt.Errorf("could not load forums: %w", err)
		}

		logger.Debug("forums", zap.Reflect("forums", forums))

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
		templates.WriteForums(c.Response(), forums)

		return nil
	})

	return e.Start(fmt.Sprintf(":%d", c.Port))
}
