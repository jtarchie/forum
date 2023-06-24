//go:build test
// +build test

package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type sqlClient struct {
	db       *sql.DB
	hostname string
}

func NewSQLClient(hostname string) (Client, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("could not initialize client: %w", err)
	}

	return &sqlClient{
		db:       db,
		hostname: hostname,
	}, nil
}

func (c *sqlClient) URL() string {
	return c.hostname
}

func (c *sqlClient) Execute(statement string, args ...interface{}) error {
	_, err := c.db.Exec(statement, args)
	if err != nil {
		return fmt.Errorf("could not execute statement: %w", err)
	}

	return nil
}

func (c *sqlClient) Query(statement string, args ...interface{}) (QueryResult, error) {
	result, err := c.db.Query(statement, args)
	if err != nil {
		return nil, fmt.Errorf("could not execute statement: %w: %w", result.Err(), err)
	}

	return result, nil
}
