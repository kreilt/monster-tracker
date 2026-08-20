package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kreilt/monster-tracker/internal/handlers"
	"github.com/kreilt/monster-tracker/internal/repository"
	"google.golang.org/api/option"
)

func main() {
	dbURL := os.Getenv("MONSTER_DB")
	if dbURL == "" {
		log.Fatal("MONSTER_DB is not set")
	}

	credPath := os.Getenv("FIREBASE_CREDENTIALS")
	if credPath == "" {
		log.Fatal("FIREBASE_CREDENTIALS is not set")
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

	opt := option.WithAuthCredentialsFile(option.ServiceAccount, credPath)
	app, err := firebase.NewApp(startCtx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	client, err := app.Auth(startCtx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v", err)
	}
	_ = client

	flavorsRepo := repository.NewFlavor(pool)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.Health(pool))
	mux.HandleFunc("GET /api/v1/flavors", handlers.Flavors(flavorsRepo))

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
