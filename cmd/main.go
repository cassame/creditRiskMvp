package main

import (
	"context"
	"credit-risk-mvp/internal/config"
	"credit-risk-mvp/internal/handlers"
	"credit-risk-mvp/internal/repository"
	"credit-risk-mvp/notifier"
	"credit-risk-mvp/services"
	"credit-risk-mvp/storage"

	"github.com/pressly/goose/v3"

	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	cfg := config.LoadConfig()
	db := storage.OpenDB(cfg)
	defer db.Close()

	log.Println("Starting migrations...")
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations finished successfully!")
	realRepo := repository.NewSqlRepository(db)
	kafkaProv := services.KafkaProducer{}
	services.InitProducer(cfg.KafkaBrokers)
	go services.StartConsumer("application_topic")

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
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("could not listen on %s: %v", ":8080", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("could not gracefully shutdown the server: %v", err)
	}
	log.Println("Server gracefully stopped")
}
