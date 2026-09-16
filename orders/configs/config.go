package configs

import (
	"os"
	"time"
)

type Config struct {
	HTTPAddress         string
	GRPCAddress         string
	MongoURI            string
	MongoDatabase       string
	PaymentsGRPCAddress string
	ShutdownTimeout     time.Duration
}

func Load() Config {
	return Config{
		HTTPAddress:         env("HTTP_ADDRESS", ":8080"),
		GRPCAddress:         env("GRPC_ADDRESS", ":9090"),
		MongoURI:            env("MONGO_URI", "mongodb://localhost:27018"),
		MongoDatabase:       env("MONGO_DATABASE", "orders_db"),
		PaymentsGRPCAddress: env("PAYMENTS_GRPC_ADDRESS", "localhost:9092"),
		ShutdownTimeout:     10 * time.Second,
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
