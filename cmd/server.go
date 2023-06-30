package cmd

import (
	"fmt"

	"github.com/jtarchie/forum/db"
	"github.com/jtarchie/forum/server"
	"github.com/jtarchie/forum/services"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"go.uber.org/zap"
)

type ServerCmd struct {
	Port     int    `help:"port to run http server" default:"8080"`
	DBServer string `help:"URL to the rqlite API" default:"http://localhost:4001"`
	Github   struct {
		Key      string `help:"key to the application" env:"GITHUB_KEY"`
		Secret   string `help:"secret to the application" env:"GITHUB_SECRET"`
		Callback string `help:"callback URL to the auth endpoint" env:"GITHUB_CALLBACK" default:"http://localhost:8080/auth/github/callback"`
	} `embed:"" prefix:"github."`
	SessionSecret string `help:"used to encrypt the session cookie" env:"SESSION_SECRET" required:""`
}

func (c *ServerCmd) Run() error {
	if c.Github.Key != "" {
		goth.UseProviders(
			github.New(c.Github.Key, c.Github.Secret, c.Github.Callback),
		)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("could not create logger: %w", err)
	}

	// run migrations
	client, err := db.NewClient(c.DBServer)
	if err != nil {
		return fmt.Errorf("could create client: %w", err)
	}

	err = services.Migration(client, logger)
	if err != nil {
		return fmt.Errorf("could not migrate: %w", err)
	}

	return server.New(
		logger,
		client,
		c.SessionSecret,
	).Start(c.Port)
}
