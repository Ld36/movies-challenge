package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName                 string
	AppEnv                  string
	HTTPPort                string
	MovieServiceGRPCAddress string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppName:                 getEnv("APP_NAME", "api-gateway"),
		AppEnv:                  getEnv("APP_ENV", "development"),
		HTTPPort:                getEnv("HTTP_PORT", "8080"),
		MovieServiceGRPCAddress: getEnv("MOVIE_SERVICE_GRPC_ADDRESS", "localhost:50051"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
