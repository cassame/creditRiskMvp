package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	KafkaBrokers     []string
	DBConnString     string
	BankruptcyURL    string
	TerroristURL     string
	HTTPtimeout      time.Duration
	DatabaseURL      string
	CreditHistoryURL string
}

func LoadConfig() Config {
	timeoutMs := 1000
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	brokersStr := getenv("KAFKA_BROKERS", "localhost:9092")
	brokers := strings.Split(brokersStr, ",")
	if v := os.Getenv("HTTP_TIMEOUT_MS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			timeoutMs = n
		}
	}
	return Config{
		KafkaBrokers:     brokers,
		BankruptcyURL:    getenv("BANKRUPTCY_URL", "http://localhost:8080/mock/bankruptcy"),
		TerroristURL:     getenv("TERRORIST_LIST_URL", "http://localhost:8080/mock/terrorist/list"),
		HTTPtimeout:      time.Duration(timeoutMs) * time.Millisecond,
		DatabaseURL:      getenv("DATABASE_URL", "host=localhost port=5432 user=user password=pass dbname=credit_risk sslmode=disable"),
		CreditHistoryURL: getenv("CREDIT_HISTORY_URL", "http://localhost:8080/mock/credit-history"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
