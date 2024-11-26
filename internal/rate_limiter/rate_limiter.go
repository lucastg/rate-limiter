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
	// Verificar se a chave está bloqueada
	if rl.config.Persistence.IsBlocked(key) {
		log.Printf("Chave %s está bloqueada", key)
		return false
	}

	// Obter o contador atual para a chave
	currentCount, err := rl.config.Persistence.GetCount(key)
	if err != nil {
		log.Printf("Erro ao acessar persistência para %s: %v", key, err)
		return false
	}

	// Incrementar o contador antes de verificar o limite
	currentCount++
	if currentCount > rl.config.Limit {
		log.Printf("Limite atingido para %s (count: %d, limit: %d)", key, currentCount, rl.config.Limit)
		// Bloquear a chave por um período definido
		if err := rl.config.Persistence.Block(key, rl.config.BlockTime); err != nil {
			log.Printf("Erro ao bloquear chave %s: %v", key, err)
		}
		return false
	}

	// Persistir o incremento no contador
	if err := rl.config.Persistence.Increment(key); err != nil {
		log.Printf("Erro ao incrementar contador para %s: %v", key, err)
		return false
	}

	// Permitir a requisição
	log.Printf("Requisição permitida para %s (count: %d)", key, currentCount)
	return true
}
