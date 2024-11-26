package persistence

import "time"

type Persistence interface {
	GetCount(key string) (int, error)
	Increment(key string) error
	Block(key string, duration time.Duration) error
	IsBlocked(key string) bool
}
