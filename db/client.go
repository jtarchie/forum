package db

import (
	"fmt"
	"strings"
)

func NewClient(hostname string) (Client, error) {
	if strings.HasPrefix(hostname, "http://") || strings.HasPrefix(hostname, "https://") {
		return NewRQLClient(hostname)
	}

	return nil, fmt.Errorf("the schema in %q is not supported", hostname)
}
