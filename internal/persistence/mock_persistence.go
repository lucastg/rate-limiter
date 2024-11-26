package persistence

import (
	"time"
)

// MockPersistence simula a persistência para fins de teste
type MockPersistence struct {
	data    map[string]int
	blocked map[string]time.Time
}

// Novo MockPersistence
func NewMockPersistence() *MockPersistence {
	return &MockPersistence{
		data: make(map[string]int),
		// blocked: make(map[string]time.Time),
	}
}

// Implementação do método GetCount simulado
func (mp *MockPersistence) GetCount(key string) (int, error) {
	count, exists := mp.data[key]
	if !exists {
		return 0, nil
	}
	return count, nil
}

// Implementação do método Increment simulado
func (mp *MockPersistence) Increment(key string) error {
	mp.data[key]++
	return nil
}

// Implementação do método Block simulado
func (mp *MockPersistence) Block(key string, duration time.Duration) error {
	if mp.blocked == nil {
		mp.blocked = make(map[string]time.Time)
	}
	mp.blocked[key] = time.Now().Add(duration)
	return nil
}

// IsBlocked verifica se a chave está bloqueada
func (mp *MockPersistence) IsBlocked(key string) bool {
	blockTime, exists := mp.blocked[key]
	if !exists {
		return false
	}

	// Verifica se o bloqueio expirou
	if time.Now().After(blockTime) {
		// Remove o bloqueio, pois o tempo expirou
		delete(mp.blocked, key)
		return false
	}

	return true
}
