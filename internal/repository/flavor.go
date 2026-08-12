package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kreilt/monster-tracker/internal/model"
)

type FlavorRepository struct {
	pool *pgxpool.Pool
}

func NewFlavorRepository(pool *pgxpool.Pool) *FlavorRepository {
	return &FlavorRepository{
		pool: pool,
	}
}

func (r *FlavorRepository) GetAll(ctx context.Context) ([]model.Flavor, error) {
	rows, err := r.pool.Query(ctx, `SELECT flavor_id, title, lineup, description, rare, region, color, status, photo FROM flavors ORDER BY flavor_id`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	flavors := []model.Flavor{}
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
