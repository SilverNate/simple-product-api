package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	ServerPort   string
	PostgresDSN  string
	RedisAddress string
}

func Load() *Config {
	_ = godotenv.Load(".env") // load .env file

	return &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		PostgresDSN:  getEnv("POSTGRES_DSN", ""),
		RedisAddress: getEnv("REDIS_ADDRESS", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	logrus.Infof("[CONFIG] ENV '%s' not found, using default: %s", key, fallback)
	return fallback
}
