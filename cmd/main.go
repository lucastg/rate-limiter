package main

import (
	"log"
	"net/http"
	"time"

	"github.com/lucastg/rate-limiter/config"
	"github.com/lucastg/rate-limiter/internal/middleware"
	"github.com/lucastg/rate-limiter/internal/persistence"
	"github.com/lucastg/rate-limiter/internal/rate_limiter"

	"github.com/gorilla/mux"
)

func main() {
	config.LoadConfig()

	redisAddr := config.GetEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := config.GetEnv("REDIS_PASSWORD", "")
	redisPersistence := persistence.NewRedisPersistence(redisAddr, redisPassword, 0)

	limiter := rate_limiter.NewRateLimiter(rate_limiter.Config{
		Limit:       10,
		BlockTime:   1 * time.Minute,
		Persistence: redisPersistence,
	})

	r := mux.NewRouter()

	r.Use(middleware.RateLimiterMiddleware(limiter))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome!"))
	}).Methods("GET")

	serverAddr := config.GetEnv("SERVER_ADDR", ":8080")
	log.Printf("Server started on %s\n", serverAddr)

	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
