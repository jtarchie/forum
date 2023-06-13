package cache

import (
	"fmt"
	"sync"
	"time"
)

type CacheFunction[T any, R any] struct {
	cached   R
	duration time.Duration
	lock     *sync.Mutex
	original func(T) (R, error)
	ttl      time.Time
}

func NewFunc[T any, R any](original func(T) (R, error), duration time.Duration) *CacheFunction[T, R] {
	return &CacheFunction[T, R]{
		duration: duration,
		lock:     &sync.Mutex{},
		original: original,
		ttl:      time.Now(),
	}
}

func (cf *CacheFunction[T, R]) Invoke(arg T) (R, error) {
	cf.lock.Lock()
	defer cf.lock.Unlock()

	var err error

	if time.Now().After(cf.ttl) {
		cf.cached, err = cf.original(arg)

		cf.ttl = time.Now().Add(cf.duration)
		if err != nil {
			return cf.cached, fmt.Errorf("could not load cached value: %w", err)
		}
	}

	return cf.cached, nil
}
