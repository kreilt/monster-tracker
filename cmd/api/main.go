package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

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

func flavorsHandler(flavorsRepo *repository.FlavorRepository) http.HandlerFunc {
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

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("unable to create pool: %v", err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}

	flavorsRepo := repository.NewFlavorRepository(pool)

	http.HandleFunc("/health", healthHandler(pool))
	http.HandleFunc("GET /api/v1/flavors", flavorsHandler(flavorsRepo))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
