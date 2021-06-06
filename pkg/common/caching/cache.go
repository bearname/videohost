package caching

import "errors"

var (
	ErrCacheUnavailable = errors.New("cache unavailable")
)

type Cache interface {
	IsOk() bool
	Get(key string) (string, error)
	Set(key string, value string) error
	Del(key string) error
}
