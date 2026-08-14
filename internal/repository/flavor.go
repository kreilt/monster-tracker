package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kreilt/monster-tracker/internal/model"
)

type Flavor struct {
	pool *pgxpool.Pool
}

func NewFlavor(pool *pgxpool.Pool) *Flavor {
	return &Flavor{
		pool: pool,
	}
}

func (r *Flavor) GetAll(ctx context.Context, lineup string) ([]model.Flavor, error) {
	rows, err := r.pool.Query(ctx, `SELECT flavor_id, title, lineup, description, rare, region, color, status, photo 
									FROM flavors 
									WHERE $1= '' OR lineup = $1
									ORDER BY flavor_id`, lineup)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	flavors := make([]model.Flavor, 0, 200)
	for rows.Next() {
		var f model.Flavor
		if err := rows.Scan(&f.FlavorID, &f.Title, &f.Lineup, &f.Description, &f.Rare, &f.Region, &f.Color, &f.Status, &f.Photo); err != nil {
			return nil, err
		}
		flavors = append(flavors, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return flavors, nil
}
