package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/jtarchie/forum/cache"
	"github.com/jtarchie/forum/db"
	"github.com/jtarchie/forum/services"
	"github.com/jtarchie/forum/templates"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Server struct {
	logger        *zap.Logger
	client        db.Client
	sessionSecret string
}

func New(
	logger *zap.Logger,
	client db.Client,
	sessionSecret string,
) *Server {
	return &Server{
		logger:        logger,
		client:        client,
		sessionSecret: sessionSecret,
	}
}

func (s *Server) Start(port int) error {
	e := echo.New()
	setupMiddleware(
		e,
		s.logger,
		s.sessionSecret,
	)

	cachedForums := cache.NewFunc(services.ListForums, time.Minute)

	e.GET("/", func(c echo.Context) error {
		forums, err := cachedForums.Invoke(s.client)
		if err != nil {
			return fmt.Errorf("could not load forums: %w", err)
		}

		email, err := SessionGet[string]("email", c)
		if err != nil {
			return fmt.Errorf("could not access session: %w", err)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
		templates.WriteLayout(c.Response(), email, templates.NewListForums(forums))

		return nil
	})

	e.GET("/forums/:name", func(c echo.Context) error {
		parts := strings.Split(c.Param("name"), "-")
		id := parts[len(parts)-1]

		return c.String(200, id)
	})

	setupAuth(e)

	return e.Start(fmt.Sprintf(":%d", port))
}
