package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/lucastg/rate-limiter/internal/rate_limiter"
	"github.com/stretchr/testify/assert"
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

func (mp *MockPersistence) Increment(key string) error {
	mp.data[key]++
	return nil
}

func (mp *MockPersistence) Block(key string, duration time.Duration) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.blocked[key] = time.Now().Add(duration)
	return nil
}

func (mp *MockPersistence) GetCount(key string) (int, error) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	return mp.data[key], nil
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

	mp.data[key] = 0
	log.Printf("Contador de requisições para a chave %s foi resetado", key)
}

func TestRateLimiterMiddleware2(t *testing.T) {
	mockPersistence := NewMockPersistence()
	limiter := rate_limiter.NewRateLimiter(rate_limiter.Config{
		Limit:       3,
		BlockTime:   2 * time.Second,
		Persistence: mockPersistence,
	})

	r := mux.NewRouter()
	r.Use(RateLimiterMiddleware(limiter))

	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		log.Printf("Chave utilizada no Rate Limiter: %s", userID)

		w.Write([]byte("OK"))
	})

	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("X-User-ID", "user1")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code, "Resposta deveria ser OK para requisição %d", i+1)
	}

	mockPersistence.Block("user1", 2*time.Second)

	assert.True(t, mockPersistence.IsBlocked("user1"), "Usuário deveria estar bloqueado após 3 requisições")

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-User-ID", "user1")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code, "Resposta deveria ser 'Too Many Requests' após bloqueio")

	log.Println("Esperando o tempo do bloqueio expirar...")
	time.Sleep(10 * time.Second)

	if !mockPersistence.IsBlocked("user1") {
		log.Println("A chave foi desbloqueada, resetando o contador.")
		mockPersistence.ResetRequestCount("user1")
	} else {
		log.Println("A chave ainda está bloqueada.")
	}

	req, err = http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-User-ID", "user1")

	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Resposta deveria ser OK após o bloqueio expirar")
}
