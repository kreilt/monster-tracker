package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kreilt/monster-tracker/internal/model"
)

type User struct {
	pool *pgxpool.Pool
}

func NewUser(pool *pgxpool.Pool) *User {
	return &User{
		pool: pool,
	}
}

func (r *User) GetOrCreate(ctx context.Context, externalID string, email string, nickname *string) (model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (external_id, email, nickname)
		VALUES ($1, $2, $3)
		ON CONFLICT (external_id) DO UPDATE SET email = EXCLUDED.email, nickname = EXCLUDED.nickname
		RETURNING user_id, external_id, nickname, email, created_at`,
		externalID, email, nickname).Scan(
		&u.UserID,
		&u.ExternalID,
		&u.Nickname,
		&u.Email,
		&u.CreatedAt)
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}
