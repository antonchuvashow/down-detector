package gin

import (
	"os"
	"strconv"

	gingonic "github.com/gin-gonic/gin"
)

func ConfigFromEnv(handlers Handlers) Config {
	return Config{
		Port:     getEnvInt("SERVER_PORT", 5436),
		Mode:     getEnv("GIN_MODE", gingonic.TestMode),
		Handlers: handlers,
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
