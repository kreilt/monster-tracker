package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func handler(pool *pgxpool.Pool) http.HandlerFunc {
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

	http.HandleFunc("/health", handler(pool))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
