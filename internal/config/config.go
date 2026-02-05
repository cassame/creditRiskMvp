package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	BankruptcyURL    string
	TerroristURL     string
	HTTPtimeout      time.Duration
	DatabaseURL      string
	CreditHistoryURL string
}

func LoadConfig() Config {
	timeoutMs := 1000
	if v := os.Getenv("HTTP_TIMEOUT_MS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			timeoutMs = n
		}
	}
	return Config{
		BankruptcyURL:    getenv("BANKRUPTCY_URL", "http://localhost:8080/mock/bankruptcy"),
		TerroristURL:     getenv("TERRORIST_LIST_URL", "http://localhost:8080/mock/terrorist/list"),
		HTTPtimeout:      time.Duration(timeoutMs) * time.Millisecond,
		DatabaseURL:      getenv("DATABASE_URL", "postgres://app:app@localhost:5432/creditrisk?sslmode=disable"),
		CreditHistoryURL: getenv("CREDIT_HISTORY_URL", "http://localhost:8080/mock/credit-history"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
