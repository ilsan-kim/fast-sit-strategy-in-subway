package cache

import (
	"time"
)

//go:generate mockgen -source=./cache.go  -destination=./mock_cache/mock_cache.go
type Cache interface {
	Get(string) (any, bool)
	Set(string, any, time.Duration)
}
