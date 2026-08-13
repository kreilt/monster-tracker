package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kreilt/monster-tracker/internal/repository"
)

func healthHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := pool.Ping(r.Context()); err != nil {
			log.Printf("unable to reach database: %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(
				map[string]string{
					"status": "error",
					"db":     "unreachable"})
			return
		}
		_ = json.NewEncoder(w).Encode(
			map[string]string{
				"status": "ok",
				"db":     "connected"})
	}
}

func flavorsHandler(flavorsRepo *repository.Flavor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		flavors, err := flavorsRepo.GetAll(r.Context())
		if err != nil {
			log.Printf("failed to get flavors, %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(
				map[string]string{
					"error": "internal"})
			return
		}
		_ = json.NewEncoder(w).Encode(flavors)
	}
}

func main() {
	dbURL := os.Getenv("MONSTER_DB")
	if dbURL == "" {
		log.Fatal("MONSTER_DB is not set")
	}

	startCtx := context.Background()

	pool, err := pgxpool.New(startCtx, dbURL)
	if err != nil {
		log.Fatalf("unable to create pool: %v", err)
	}

	defer pool.Close()

	if err := pool.Ping(startCtx); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}

	flavorsRepo := repository.NewFlavor(pool)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler(pool))
	mux.HandleFunc("GET /api/v1/flavors", flavorsHandler(flavorsRepo))

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	log.Println("server started on :8080")

	<-ctx.Done()
	log.Println("shutdown")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
