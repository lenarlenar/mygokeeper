package server

import (
	"os"
)

// Config содержит конфигурацию сервера.
type Config struct {
	ServerAddr     string
	JWTSecret      []byte
	DBConn         string
	MigrationsPath string
}

// Load загружает конфигурацию из переменных окружения с дефолтами.
func Load() *Config {
	addr := os.Getenv("GOKEEPER_SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	jwt := os.Getenv("GOKEEPER_JWT_SECRET")
	if jwt == "" {
		jwt = "my-secret-key"
	}

	db := os.Getenv("GOKEEPER_DB")
	if db == "" {
		db = "postgres://postgres:postgres@localhost:5432/gokeeper?sslmode=disable"
	}

	migrations := os.Getenv("GOKEEPER_MIGRATIONS")
	if migrations == "" {
		migrations = "../../internal/server/migrations"
	}

	return &Config{
		ServerAddr:     addr,
		JWTSecret:      []byte(jwt),
		DBConn:         db,
		MigrationsPath: migrations,
	}
}
