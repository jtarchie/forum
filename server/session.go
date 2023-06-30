package server

import (
	"fmt"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const sessionName = "session"

func SessionPut[T any](name string, value T, c echo.Context) error {
	state, err := session.Get(sessionName, c)
	if err != nil {
		return fmt.Errorf("could not find session named %q: %w", sessionName, err)
	}

	state.Values[name] = value

	err = state.Save(c.Request(), c.Response())
	if err != nil {
		return fmt.Errorf("could not save session: %w", err)
	}

	return nil
}

func SessionGet[T any](name string, c echo.Context) (T, error) {
	var zeroValue T

	state, err := session.Get(sessionName, c)
	if err != nil {
		return zeroValue, fmt.Errorf("could not find session named %q: %w", sessionName, err)
	}

	value, ok := state.Values[name]
	if !ok {
		return zeroValue, nil
	}

	return value.(T), nil
}
