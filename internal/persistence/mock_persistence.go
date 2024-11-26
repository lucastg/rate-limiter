package persistence

import (
	"log"
	"sync"
	"time"
)

type MockPersistence struct {
	mu      sync.Mutex
	data    map[string]int
	blocked map[string]time.Time
}

func NewMockPersistence() *MockPersistence {
	return &MockPersistence{
		data:    make(map[string]int),
		blocked: make(map[string]time.Time),
	}
}

func (mp *MockPersistence) GetCount(key string) (int, error) {
	count, exists := mp.data[key]
	if !exists {
		return 0, nil
	}
	return count, nil
}

func (mp *MockPersistence) Increment(key string) error {
	mp.data[key]++
	return nil
}

func (mp *MockPersistence) Block(key string, duration time.Duration) error {
	if mp.blocked == nil {
		mp.blocked = make(map[string]time.Time)
	}

	blockUntil := time.Now().Add(duration)
	mp.blocked[key] = blockUntil
	log.Printf("Chave %s foi bloqueada até %s", key, blockUntil)
	return nil
}

func (mp *MockPersistence) IsBlocked(key string) bool {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	blockTime, exists := mp.blocked[key]
	if !exists {
		return false
	}

	if time.Now().After(blockTime) {
		delete(mp.blocked, key)
		return false
	}

	return true
}

func (mp *MockPersistence) ResetRequestCount(key string) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if _, exists := mp.blocked[key]; !exists {
		mp.data[key] = 0
		log.Printf("Contador de requisições para a chave %s foi resetado", key)
	} else {
		log.Printf("A chave %s ainda está bloqueada, não é possível resetar o contador", key)
	}
}
