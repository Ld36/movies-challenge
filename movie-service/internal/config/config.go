package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName         string
	AppEnv          string
	GRPCPort        string
	MongoURI        string
	MongoDatabase   string
	MongoCollection string
	RabbitMQURL     string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppName:         getEnv("APP_NAME", "movie-service"),
		AppEnv:          getEnv("APP_ENV", "development"),
		GRPCPort:        getEnv("GRPC_PORT", "50051"),
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:   getEnv("MONGO_DATABASE", "movies_db"),
		MongoCollection: getEnv("MONGO_COLLECTION", "movies"),
		RabbitMQURL:     getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
