package persistence

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisPersistence(t *testing.T) {
	mockPersistence := &MockPersistence{
		data: make(map[string]int),
	}

	// Teste de incremento
	err := mockPersistence.Increment("user1")
	assert.NoError(t, err)
	count, err := mockPersistence.GetCount("user1")
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Teste de bloqueio (não implementado no mock, mas deveria verificar o estado de bloqueio)
	mockPersistence.Block("user1", 1*time.Minute)
	// Adicione sua lógica de verificação de bloqueio aqui

	// Testando a contagem após o incremento
	err = mockPersistence.Increment("user1")
	assert.NoError(t, err)
	count, err = mockPersistence.GetCount("user1")
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}
