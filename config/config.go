package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	envFile := "../cmd/.env"
	if err := godotenv.Load(envFile); err != nil {
		log.Println("Warning: No .env file found at", envFile)
	}
}

func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func GetEnvAsInt(key string, fallback int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Unable to parse %s as int, using fallback %d", key, fallback)
		return fallback
	}
	return value
}
