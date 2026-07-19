package repository

import (
	"context"
	"rest-api/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateBlog(pool *pgxpool.Pool, title string, content string)(*models.Blog, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	
	ctx, cancel = context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO blogs (title, content, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, title, content, created_at, updated_at
	`
	var blog models.Blog
	var err error = pool.QueryRow(ctx,query,title,content).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &blog, nil
}