package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/kreilt/monster-tracker/internal/model"
	"github.com/kreilt/monster-tracker/internal/repository"
)

type contextKey string

const userKey contextKey = "user"

func UserFromContext(ctx context.Context) (model.User, bool) {
	u, ok := ctx.Value(userKey).(model.User)
	return u, ok
}

func Auth(client *auth.Client, userRepo *repository.User) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				writeError(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			idToken := strings.TrimPrefix(header, "Bearer ")
			if idToken == "" {
				writeError(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			token, err := client.VerifyIDToken(r.Context(), idToken)
			if err != nil {
				writeError(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			email, _ := token.Claims["email"].(string)
			name, _ := token.Claims["name"].(string)
			var nickname *string
			if name != "" {
				nickname = &name
			}

			user, err := userRepo.GetOrCreate(r.Context(), token.UID, email, nickname)
			if err != nil {
				log.Printf("failed to get or create user: %v", err)
				writeError(w, "internal", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), userKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})

	}
}
