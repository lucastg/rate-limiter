package rate_limiter

import (
	"testing"
	"time"

	"github.com/lucastg/rate-limiter/internal/persistence"
	"github.com/stretchr/testify/assert"
)

func TestAllowBeforeLimit(t *testing.T) {
	mockPersistence := persistence.NewMockPersistence()

	limiter := NewRateLimiter(Config{
		Limit:       5,
		BlockTime:   1 * time.Minute,
		Persistence: mockPersistence,
	})

	for i := 0; i < 5; i++ {
		allowed := limiter.Allow("user1")
		assert.True(t, allowed, "Deveria ser permitido")
	}

	allowed := limiter.Allow("user1")
	assert.False(t, allowed, "Não deveria ser permitido após atingir o limite")
}

func TestBlockAfterLimit(t *testing.T) {
	mockPersistence := persistence.NewMockPersistence()

	limiter := NewRateLimiter(Config{
		Limit:       5,
		BlockTime:   3 * time.Second,
		Persistence: mockPersistence,
	})

	for i := 0; i < 5; i++ {
		allowed := limiter.Allow("user1")
		assert.True(t, allowed, "Deveria ser permitido antes de atingir o limite")
	}

	blocked := limiter.Allow("user1")
	assert.False(t, blocked, "O usuário deveria estar bloqueado após atingir o limite")

	time.Sleep(5 * time.Second)

	allowed := limiter.Allow("user1")
	assert.False(t, allowed, "Deveria ser permitido após o tempo de bloqueio expirar")
}
