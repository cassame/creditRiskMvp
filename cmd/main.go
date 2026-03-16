package main

import (
	"context"
	"credit-risk-mvp/internal/config"
	"credit-risk-mvp/internal/handlers"
	"credit-risk-mvp/internal/logger"
	"credit-risk-mvp/internal/repository"
	"credit-risk-mvp/notifier"
	"credit-risk-mvp/services"
	"credit-risk-mvp/storage"

	"github.com/pressly/goose/v3"

	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	cfg := config.LoadConfig()
	logger.InitLogger("local")
	logger.Lg.Info("Starting application", "port", "8080", "env", "local")
	db := storage.OpenDB(cfg)
	defer db.Close()

	logger.Lg.Info("Starting migrations")
	if err := goose.Up(db, "migrations"); err != nil {
		logger.Lg.Error("Migration failed", "error", err)
		os.Exit(1)
	}
	logger.Lg.Info("Migrations finished successfully")

	realRepo := repository.NewSqlRepository(db)
	kafkaProv := services.KafkaProducer{}
	services.InitProducer(cfg.KafkaBrokers)
	logger.Lg.Info("Starting Kafka consumer", "topic", "application_topic")
	go services.StartConsumer(cfg.KafkaBrokers, "application_topic")

	notifier := notifier.LogNotifier{}
	appHandler := handlers.NewApplicationsHandler(cfg, realRepo, kafkaProv, notifier)

	mux := http.NewServeMux()
	mux.Handle("/applications", appHandler)
	mux.HandleFunc("/applications/", handlers.MakeGetApplicationHandler(realRepo))
	mux.HandleFunc("/mock/bankruptcy", handlers.HandleMockBankruptcy)
	mux.HandleFunc("/mock/terrorist/list", handlers.HandleMockTerroristList)
	mux.HandleFunc("/mock/credit-history", handlers.HandleMockCreditHistory)
	srv := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		logger.Lg.Info("Server is listening", "addr", ":8080")
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.Lg.Error("Server crash", "error", err)
			os.Exit(1)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Lg.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Lg.Error("Graceful shutdown failed", "error", err)
	}
	logger.Lg.Info("Server gracefully stopped")
}
