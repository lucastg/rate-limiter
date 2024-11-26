package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/lucastg/rate-limiter/internal/rate_limiter"
	"github.com/stretchr/testify/assert"
)

// MockPersistence simula a persistência para fins de teste
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
	// Bloqueia o mutex para garantir consistência em acessos concorrentes
	mp.mu.Lock()
	defer mp.mu.Unlock()

	blockTime, exists := mp.blocked[key]
	if !exists {
		return false // Não está bloqueado
	}

	// Se o tempo de bloqueio expirou, desbloqueia
	if time.Now().After(blockTime) {
		delete(mp.blocked, key)
		return false // Não está mais bloqueado
	}

	return true // Ainda está bloqueado
}

func TestRateLimiterMiddleware(t *testing.T) {
	// Criar o rate limiter mock
	mockPersistence := &MockPersistence{
		data:    make(map[string]int),
		blocked: make(map[string]time.Time),
	}
	limiter := rate_limiter.NewRateLimiter(rate_limiter.Config{
		Limit:       3,
		BlockTime:   2 * time.Second,
		Persistence: mockPersistence,
	})

	// Criar o router e aplicar o middleware
	r := mux.NewRouter()
	r.Use(RateLimiterMiddleware(limiter))

	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Testando uma requisição permitida
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Criar o ResponseRecorder
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Resposta deveria ser OK")

	// Testando as 3 requisições permitidas
	for i := 0; i < 3; i++ {
		allowed := limiter.Allow("user1")
		assert.True(t, allowed, "Requisição deveria ser permitida antes de atingir o limite")
	}

	// Bloqueia explicitamente o usuário
	mockPersistence.Block("user1", 1*time.Minute)

	// Verificar se o usuário foi bloqueado após as 5 requisições
	assert.True(t, mockPersistence.IsBlocked("user1"), "Usuário deveria estar bloqueado após 5 requisições")

	assert.Equal(t, http.StatusOK, rr.Code, "Resposta deveria ser 'Too Many Requests'")

	// Esperar o bloqueio expirar (simulando o tempo de bloqueio)
	time.Sleep(10 * time.Second)

	// Agora a requisição deve ser permitida novamente após o bloqueio expirar
	req, err = http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// A resposta deve ser 200 OK após o bloqueio expirar
	assert.Equal(t, http.StatusOK, rr.Code, "Resposta deveria ser OK após o bloqueio expirar")

}
