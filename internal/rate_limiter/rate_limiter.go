package rate_limiter

import (
	"log"
	"time"

	"github.com/lucastg/rate-limiter/internal/persistence"
)

type RateLimiter interface {
	Allow(key string) bool
}

type Config struct {
	Limit       int
	BlockTime   time.Duration
	Persistence persistence.Persistence
}

type rateLimiter struct {
	config Config
}

func NewRateLimiter(config Config) RateLimiter {
	return &rateLimiter{config: config}
}

func (rl *rateLimiter) Allow(key string) bool {
	if rl.config.Persistence.IsBlocked(key) {
		log.Printf("Chave %s está bloqueada", key)
		return false
	}

	currentCount, err := rl.config.Persistence.GetCount(key)
	if err != nil {
		log.Printf("Erro ao acessar persistência para %s: %v", key, err)
		return false
	}

	currentCount++
	if currentCount > rl.config.Limit {
		log.Printf("Limite atingido para %s (count: %d, limit: %d)", key, currentCount, rl.config.Limit)
		if err := rl.config.Persistence.Block(key, rl.config.BlockTime); err != nil {
			log.Printf("Erro ao bloquear chave %s: %v", key, err)
		}
		return false
	}

	if err := rl.config.Persistence.Increment(key); err != nil {
		log.Printf("Erro ao incrementar contador para %s: %v", key, err)
		return false
	}

	log.Printf("Requisição permitida para %s (count: %d)", key, currentCount)
	return true
}
