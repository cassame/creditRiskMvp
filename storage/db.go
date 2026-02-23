package storage

import (
	"credit-risk-mvp/internal/config"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func OpenDB(cfg config.Config) *sql.DB {
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("cannot open storage:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("cannot ping storage:", err)
	}
	return db
}
