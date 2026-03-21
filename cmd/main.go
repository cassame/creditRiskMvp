package main

import (
	"context"
	"credit-risk-mvp/internal/config"
	"credit-risk-mvp/internal/domain"
	"credit-risk-mvp/internal/handlers"
	"credit-risk-mvp/internal/logger"
	"credit-risk-mvp/internal/repository"
	"credit-risk-mvp/notifier"
	"credit-risk-mvp/services"
	"credit-risk-mvp/storage"
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"

	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	cfg := config.LoadConfig()
	logger.InitLogger("local")

	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger.Lg.Info("Starting application", "port", "8080", "env", "local")

	db := storage.OpenDB(cfg)
	defer func() { _ = db.Close() }()
	runMigrations(db)

	services.InitProducer(cfg.KafkaBrokers)
	defer services.CloseProducer()

	logger.Lg.Info("Starting Kafka consumer", "topic", "application_topic")
	go services.StartConsumer(appCtx, cfg.KafkaBrokers, "application_topic")

	terrStore := initTerroristStore(appCtx, cfg)

	realRepo := repository.NewSqlRepository(db)
	kafkaProv := services.KafkaProducer{}
	notifier := notifier.LogNotifier{}

	appHandler := handlers.NewApplicationsHandler(cfg, realRepo, kafkaProv, notifier, terrStore)

	srv := setupServer(appHandler, realRepo)

	go func() {
		logger.Lg.Info("Server is listening", "addr", ":8080")
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.Lg.Error("Server crash", "error", err)
			os.Exit(1)
		}
	}()
	<-appCtx.Done()
	logger.Lg.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Lg.Error("Graceful shutdown failed", "error", err)
	}
	logger.Lg.Info("Server gracefully stopped")
}

func setupServer(appHandler *handlers.ApplicationHandler, realRepo domain.Repository) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/applications", appHandler)
	mux.HandleFunc("/applications/", handlers.MakeGetApplicationHandler(realRepo))
	mux.HandleFunc("/mock/bankruptcy", handlers.HandleMockBankruptcy)
	mux.HandleFunc("/mock/terrorist/list", handlers.HandleMockTerroristList)
	mux.HandleFunc("/mock/credit-history", handlers.HandleMockCreditHistory)

	return &http.Server{Addr: ":8080", Handler: mux}
}

func initTerroristStore(ctx context.Context, cfg config.Config) domain.TerroristStore {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := rdb.Ping(pingCtx).Err(); err != nil {
		logger.Lg.Warn("failed to connect to redis, falling back to local store", "error", err)
		cfg.UseRedis = false
	}
	if cfg.UseRedis {
		logger.Lg.Info("Using Redis for terrorist store")
		return repository.NewRedisTerroristStore(rdb)
	} else {
		logger.Lg.Info("Using Local memory for terrorist store")
		return repository.NewLocalTerroristStore()
	}
}

func runMigrations(db *sql.DB) {
	logger.Lg.Info("Starting migrations")
	if err := goose.Up(db, "migrations"); err != nil {
		logger.Lg.Error("Migration failed", "error", err)
		os.Exit(1)
	}
	logger.Lg.Info("Migrations finished successfully")
}
