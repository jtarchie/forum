package db

import (
	"fmt"

	"github.com/rqlite/gorqlite"
)

type Client struct {
	conn *gorqlite.Connection
}

func NewClient(hostname string) (*Client, error) {
	conn, err := gorqlite.Open(hostname)
	if err != nil {
		return nil, fmt.Errorf("could not initialize client: %w", err)
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Execute(statement string, args ...interface{}) error {
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

type QueryResult interface {
	Next() bool
	NumRows() int64
	Scan(dest ...interface{}) error
}

func (c *Client) Query(statement string, args ...interface{}) (QueryResult, error) {
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
