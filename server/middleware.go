package server

import (
	"github.com/gorilla/sessions"
	customMiddleware "github.com/jtarchie/middleware"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func setupMiddleware(
	e *echo.Echo,
	logger *zap.Logger,
	sessionSecret string,
) {
	e.Use(middleware.Secure())
	e.Use(middleware.CSRF())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionSecret))))
	e.Use(customMiddleware.ZapLogger(logger))
}
