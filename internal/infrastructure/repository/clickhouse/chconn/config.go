package chconn

import (
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type Config struct {
	Addr     string
	Database string
	Username string
	Password string
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
