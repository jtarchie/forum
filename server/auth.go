package server

import (
	"fmt"
	"net/http"

	"github.com/jtarchie/forum/gothic"
	"github.com/labstack/echo/v4"
)

func setupAuth(e *echo.Echo) {
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

		err := SessionPut("email", "", c)
		if err != nil {
			return fmt.Errorf("could not clear user: %w", err)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/")
	})
}
