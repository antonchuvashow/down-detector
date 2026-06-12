package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConfigFromEnv() Config {
	return Config{
		Addr:     getEnv("POSTGRES_ADDR", "localhost:5432"),
		Database: getEnv("POSTGRES_DB", "detector"),
		Username: getEnv("POSTGRES_USER", "detector"),
		Password: getEnv("POSTGRES_PASSWORD", "detector"),
		SSLMode:  getEnv("POSTGRES_SSL_MODE", "disable"),

		MaxOpenConns:    25,
		MaxIdleConns:    25,
		ConnMaxLifetime: 5 * time.Minute,
		PingTimeout:     5 * time.Second,
	}
}

type Config struct {
	Addr     string
	Database string
	Username string
	Password string
	SSLMode  string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	PingTimeout     time.Duration
}

func New(cfg Config) (*sql.DB, error) {
	dsn := buildDSN(cfg)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	applyPoolConfig(db, cfg)

	pingTimeout := cfg.PingTimeout
	if pingTimeout == 0 {
		pingTimeout = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres connection: %w", err)
	}

	return db, nil
}

func buildDSN(cfg Config) string {
	sslMode := cfg.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Host:   cfg.Addr,
		Path:   "/" + cfg.Database,
	}

	q := u.Query()
	q.Set("sslmode", sslMode)
	u.RawQuery = q.Encode()

	return u.String()
}

func applyPoolConfig(db *sql.DB, cfg Config) {
	maxOpenConns := cfg.MaxOpenConns
	if maxOpenConns == 0 {
		maxOpenConns = 25
	}

	maxIdleConns := cfg.MaxIdleConns
	if maxIdleConns == 0 {
		maxIdleConns = 25
	}

	connMaxLifetime := cfg.ConnMaxLifetime
	if connMaxLifetime == 0 {
		connMaxLifetime = 5 * time.Minute
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
