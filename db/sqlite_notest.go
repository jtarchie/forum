//go:build !test
// +build !test

package db

import (
	"fmt"
)

func NewSQLClient(hostname string) (Client, error) {
	return nil, fmt.Errorf("sqlite not supported")
}
