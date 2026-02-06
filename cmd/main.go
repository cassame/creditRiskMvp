package main

import (
	"context"
	"credit-risk-mvp/internal/config"
	"credit-risk-mvp/internal/handlers"
	"credit-risk-mvp/services"
	"credit-risk-mvp/storage"

	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	services.InitProducer()

	go services.StartConsumer("application_topic")

	cfg := config.LoadConfig()
	db := storage.OpenDB(cfg)
	defer func() {
		_ = db.Close()
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/applications", handlers.MakeApplicationsHandler(cfg, db))
	mux.HandleFunc("/applications/", handlers.MakeGetApplicationHandler(db))

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
