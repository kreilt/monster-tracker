package handlers

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Health(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := pool.Ping(r.Context()); err != nil {
			log.Printf("unable to reach database: %v", err)
			writeJSON(w, map[string]string{
				"status": "error",
				"db":     "unreachable"}, http.StatusServiceUnavailable)
			return
		}
		writeJSON(w, map[string]string{
			"status": "ok",
			"db":     "connected"}, http.StatusOK)
	}
}
