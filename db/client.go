package db

import (
	"fmt"

	"github.com/rqlite/gorqlite"
)

type client struct {
	conn     *gorqlite.Connection
	hostname string
}

type QueryResult interface {
	Next() bool
	Scan(...interface{}) error
	NumRows() int64
}

type Client interface {
	Execute(string, ...interface{}) error
	Query(string, ...interface{}) (QueryResult, error)
	URL() string
}

func NewClient(hostname string) (Client, error) {
	conn, err := gorqlite.Open(hostname)
	if err != nil {
		return nil, fmt.Errorf("could not initialize client: %w", err)
	}

	return &client{
		conn:     conn,
		hostname: hostname,
	}, nil
}

func (c *client) URL() string {
	return c.hostname
}

func (c *client) Execute(statement string, args ...interface{}) error {
	if len(args) == 0 {
		rows, err := c.conn.WriteOne(statement)
		if err != nil {
			return fmt.Errorf("could not execute statement: %w: %w", rows.Err, err)
		}

		return nil
	}

	rows, err := c.conn.WriteOneParameterized(gorqlite.ParameterizedStatement{
		Query:     statement,
		Arguments: args,
	})
	if err != nil {
		return fmt.Errorf("could not execute statement: %w: %w", rows.Err, err)
	}

	return nil
}

func (c *client) Query(statement string, args ...interface{}) (QueryResult, error) {
	if len(args) == 0 {
		rows, err := c.conn.QueryOne(statement)
		if err != nil {
			return nil, fmt.Errorf("could not query statement: %w: %w", rows.Err, err)
		}

		return &rows, nil
	}

	rows, err := c.conn.QueryOneParameterized(gorqlite.ParameterizedStatement{
		Query:     statement,
		Arguments: args,
	})
	if err != nil {
		return nil, fmt.Errorf("could not query prepared statement: %w: %w", rows.Err, err)
	}

	return &rows, nil
}
