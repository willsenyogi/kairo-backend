package database

import (
	"context"
	"fmt"
	"log"

	"kairo-backend/internal/config"

	"github.com/jackc/pgx/v5"
)

func NewDB(cfg *config.Config) *pgx.Conn {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	return conn
}
