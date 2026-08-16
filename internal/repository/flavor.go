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

func (r *Flavor) List(ctx context.Context, filters FlavorFilter) ([]model.Flavor, error) {
	rows, err := r.pool.Query(ctx, `SELECT flavor_id, title, lineup, description, rare, region, color, status, photo 
									FROM flavors 
									WHERE ($1 = '' OR title ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
									AND ($2 = '' OR lineup = $2)
									AND ($3 = '' OR rare::text = $3)
									AND ($4 = '' OR region = $4)
									AND ($5 = '' OR status::text = $5)
									ORDER BY flavor_id`, filters.Search, filters.Lineup, filters.Rare, filters.Region, filters.Status)
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
