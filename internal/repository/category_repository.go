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

func GetAllCategory(pool *pgxpool.Pool)([]models.Category,error){
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	var query string =
	`
		SELECT id,name,slug,created_at,updated_at
		FROM categories
		ORDER BY created_at DESC
	`

	var rows,err = pool.Query(ctx,query)
	if err != nil{
		return nil,err
	}

	defer rows.Close()

	var categories []models.Category = []models.Category{}

	for rows.Next(){
		var category models.Category

		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil,err
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil{
		return nil,err
	}

	return categories,nil

}


func GetCategoryByID(pool *pgxpool.Pool,id int)(*models.Category,error){
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(),5 *time.Second)
	defer cancel()

	var query string =
	`
		SELECT id,name,slug,created_at,updated_at
		FROM categories
		WHERE id = $1
	`

	var categories models.Category

	err := pool.QueryRow(ctx,query,id).Scan(
		&categories.ID,
		&categories.Name,
		&categories.Slug,
		&categories.CreatedAt,
		&categories.UpdatedAt,
	)

	if err != nil{
		return nil,err
	}

	return &categories,nil
}

func UpdateCategory(pool *pgxpool.Pool,id int,name string,slug string)(*models.Category,error){
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(),5 *time.Second)
	defer cancel()

	var query string = `
	UPDATE categories
	SET name = COALESCE($1, name),
		slug = COALESCE($2, slug),
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $3 
	RETURNING id, name, slug, created_at, updated_at
	`

	var category models.Category

	var err error = pool.QueryRow(ctx,query,name,slug,id).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil{
		return nil,err
	}

	return &category,nil
}

func DeleteCategory(pool *pgxpool.Pool,id int)(*models.Category,error){
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(),5 *time.Second)
	defer cancel()

	var query string = `
	DELETE FROM categories
	WHERE ID = $1 
	RETURNING id, name, slug, created_at, updated_at
	`

	var category models.Category

	var err error = pool.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil{
		return nil,err
	}

	return &category,nil
}