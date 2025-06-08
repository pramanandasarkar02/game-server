package config

import (
	"os"
	"strconv"
	"github.com/joho/godotenv"
	"github.com/pramanandasarkar02/game-server/pkg/logger"
)

type Config struct {
	Port           string
	WSPingInterval int
	WSReadTimeout  int
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found: %v", err)
	}

	cfg := &Config{
		Port:           getEnv("PORT", "4000"),
		WSPingInterval: getEnvAsInt("WS_PING_INTERVAL", 30),
		WSReadTimeout:  getEnvAsInt("WS_READ_TIMEOUT", 60),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}