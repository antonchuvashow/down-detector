package chconn

import (
	"fmt"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type Config struct {
	Addr     string
	Database string
	Username string
	Password string
}

func ConfigFromEnv() Config {
	return Config{
		Addr:     getEnv("CLICKHOUSE_ADDR", "localhost:9000"),
		Database: getEnv("CLICKHOUSE_DB", "analytics"),
		Username: getEnv("CLICKHOUSE_USER", "dev"),
		Password: getEnv("CLICKHOUSE_PASSWORD", "password"),
	}
}

func New(cfg Config) (clickhouse.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.Addr},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("open clickhouse connection: %w", err)
	}

	return conn, nil
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
