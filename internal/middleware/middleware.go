package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/lucastg/rate-limiter/internal/rate_limiter"
)

// RateLimiterMiddleware aplica o rate limiting nas requisições.
func RateLimiterMiddleware(limiter rate_limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = strings.Split(forwarded, ",")[0] // Pega o primeiro IP (cliente original)
			}

			key := ip
			if token := r.Header.Get("API_KEY"); token != "" {
				key = token
			}

			log.Printf("Chave utilizada no Rate Limiter: %s", key)
			if !limiter.Allow(key) {
				log.Printf("Limite atingido para a chave: %s", key)
				http.Error(w, "Você atingiu o número máximo de solicitações ou ações permitidas!", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}