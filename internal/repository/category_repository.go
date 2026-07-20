package repository

import (
	"context"
	"rest-api/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateCategory(pool *pgxpool.Pool,name string,slug string)(*models.Category,error){
	var ctx context.Context
	var cancel context.CancelFunc
	
	ctx, cancel = context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	var query string =
	`
		INSERT INTO categories (name, slug) VALUES ($1, $2) 
		RETURNING id, name, slug, created_at, updated_at
	`
	var category models.Category
	var err error = pool.QueryRow(ctx,query,name,slug).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil,err
	}

	return &category,nil

}