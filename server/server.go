package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jtarchie/forum/cache"
	"github.com/jtarchie/forum/db"
	"github.com/jtarchie/forum/gothic"
	"github.com/jtarchie/forum/services"
	"github.com/jtarchie/forum/templates"
	customMiddleware "github.com/jtarchie/middleware"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.Use(middleware.Secure())
	e.Use(middleware.CSRF())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(s.sessionSecret))))
	e.Use(customMiddleware.ZapLogger(s.logger))

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

	e.GET("/auth/:provider/login", func(c echo.Context) error {
		// try to get the user without re-authenticating
		user, err := gothic.CompleteUserAuth(c)
		if err == nil {
			err = SessionPut("email", user.Email, c)
			if err != nil {
				return fmt.Errorf("could not persist user: %w", err)
			}

			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		return gothic.BeginAuthHandler(c)
	})

	e.GET("/auth/:provider/callback", func(c echo.Context) error {
		user, err := gothic.CompleteUserAuth(c)
		if err != nil {
			return fmt.Errorf("could not complete auth: %w", err)
		}

		err = SessionPut("email", user.Email, c)
		if err != nil {
			return fmt.Errorf("could not persist user: %w", err)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	e.GET("/auth/:provider/logout", func(c echo.Context) error {
		_ = gothic.Logout(c)

		sess, err := session.Get("session", c)
		if err != nil {
			return fmt.Errorf("could not load session: %w", err)
		}

		sess.Values["email"] = ""
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			return fmt.Errorf("could not save session: %w", err)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	return e.Start(fmt.Sprintf(":%d", port))
}
