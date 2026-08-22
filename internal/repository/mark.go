package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kreilt/monster-tracker/internal/model"
)

type Mark struct {
	pool *pgxpool.Pool
}

func NewMark(pool *pgxpool.Pool) *Mark {
	return &Mark{
		pool: pool,
	}
}

func (r *Mark) Upsert(ctx context.Context, m model.Mark) (model.Mark, error) {
	var out model.Mark
	err := r.pool.QueryRow(ctx,
		`INSERT INTO flavors_users (user_id, flavor_id, status_tried, tried_at, user_photo)
		VALUES($1, $2, $3, CASE WHEN $3 THEN now() ELSE NULL END, $4)
		ON CONFLICT (user_id, flavor_id) DO UPDATE SET 
			status_tried = EXCLUDED.status_tried, 	
			tried_at = CASE WHEN EXCLUDED.status_tried 
							THEN COALESCE(flavors_users.tried_at, now()) 
							ELSE NULL END,
			user_photo = COALESCE(EXCLUDED.user_photo, flavors_users.user_photo)
		RETURNING flavor_user_id, user_id, flavor_id, status_tried, tried_at, user_photo`,
		m.UserID, m.FlavorID, m.StatusTried, m.UserPhoto).Scan(
		&out.MarkID,
		&out.UserID,
		&out.FlavorID,
		&out.StatusTried,
		&out.TriedAt,
		&out.UserPhoto,
	)
	if err != nil {
		return model.Mark{}, err
	}
	return out, nil
}

func (r *Mark) Delete(ctx context.Context, userID, flavorID int) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM flavors_users WHERE user_id = $1 AND flavor_id = $2`,
		userID, flavorID,
	)
	return err
}

func (r *Mark) ListByUser(ctx context.Context, userID int) ([]model.Mark, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT flavor_user_id, user_id, flavor_id, status_tried, tried_at, user_photo
		FROM flavors_users
		WHERE user_id = $1
		ORDER BY flavor_id`,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	marks := make([]model.Mark, 0)
	for rows.Next() {
		var m model.Mark
		if err := rows.Scan(&m.MarkID, &m.UserID, &m.FlavorID, &m.StatusTried, &m.TriedAt, &m.UserPhoto); err != nil {
			return nil, err
		}
		marks = append(marks, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return marks, nil
}
