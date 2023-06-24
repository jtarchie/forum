package db

import (
	"fmt"
	"strings"
)

type QueryResult interface {
	Next() bool
	Scan(...interface{}) error
}

type Client interface {
	Execute(string, ...interface{}) error
	Query(string, ...interface{}) (QueryResult, error)
	URL() string
}

func NewClient(hostname string) (Client, error) {
	if strings.HasPrefix(hostname, "http://") || strings.HasPrefix(hostname, "https://") {
		return NewRQLClient(hostname)
	}

	if strings.HasPrefix(hostname, "sqlite://") {
		return NewSQLClient(hostname)
	}

	return nil, fmt.Errorf("the schema in %q is not supported", hostname)
}
