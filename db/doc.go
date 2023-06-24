package db

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
