package rate_limiter

import (
	"testing"
	"time"

	"github.com/lucastg/rate-limiter/internal/persistence" // Certifique-se de importar o pacote correto
	"github.com/stretchr/testify/assert"
)

func TestAllowBeforeLimit(t *testing.T) {
	// Criando a persistência mock
	mockPersistence := persistence.NewMockPersistence()

	// Criando o rate limiter com o mock de persistência
	limiter := NewRateLimiter(Config{
		Limit:       5,
		BlockTime:   1 * time.Minute,
		Persistence: mockPersistence,
	})

	// Testando as 5 requisições permitidas
	for i := 0; i < 5; i++ {
		allowed := limiter.Allow("user1")
		assert.True(t, allowed, "Deveria ser permitido")
	}

	// A 6ª requisição deve ser bloqueada
	allowed := limiter.Allow("user1")
	assert.False(t, allowed, "Não deveria ser permitido após atingir o limite")
}

func TestBlockAfterLimit(t *testing.T) {
	// Criando a persistência mock
	mockPersistence := persistence.NewMockPersistence()

	// Criando o rate limiter com o mock de persistência
	limiter := NewRateLimiter(Config{
		Limit:       5,
		BlockTime:   3 * time.Second, // Bloqueio de 3 segundos
		Persistence: mockPersistence,
	})

	// Realizando 5 requisições para atingir o limite
	for i := 0; i < 5; i++ {
		allowed := limiter.Allow("user1")
		assert.True(t, allowed, "Deveria ser permitido antes de atingir o limite")
	}

	// A próxima requisição deve ser bloqueada
	blocked := limiter.Allow("user1")
	assert.False(t, blocked, "O usuário deveria estar bloqueado após atingir o limite")

	// Esperar o tempo de bloqueio expirar
	time.Sleep(5 * time.Second) // Tempo de espera maior que o tempo de bloqueio

	// Após o bloqueio expirar, a requisição deve ser permitida novamente
	allowed := limiter.Allow("user1")
	assert.False(t, allowed, "Deveria ser permitido após o tempo de bloqueio expirar")
}
