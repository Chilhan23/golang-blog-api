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

func GetAllBlogs(pool *pgxpool.Pool)([]models.Blog,error){
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	var query = `
		SELECT id,title,content,created_at,updated_at
		FROM blogs
		ORDER BY created_at DESC
	`
	var rows, err = pool.Query(ctx,query)
	if err != nil{
		return nil,err
	}

	defer rows.Close()

	var blogs []models.Blog = []models.Blog{}

	for rows.Next(){
		var blog models.Blog

		err = rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Content,
			&blog.CreatedAt,
			&blog.UpdatedAt,
		)

		if err != nil{
			return nil,err
		}

		blogs = append(blogs, blog)
	}

	if err = rows.Err(); err != nil{
		return nil,err
	}

	return blogs, nil
}

func GetBlogByID(pool *pgxpool.Pool, id int)(*models.Blog,error){
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(),5 *time.Second)
	defer cancel()

	var query = `
		SELECT id,title,content,created_at,updated_at
		FROM blogs
		WHERE ID = $1
	`
	var blog models.Blog

	err := pool.QueryRow(ctx,query,id).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)

	if err != nil{
		return nil, err
	}

	return &blog,nil

}

func UpdateBlog(pool *pgxpool.Pool,id int,title string,content string)(*models.Blog,error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(),5 *time.Second)
	defer cancel()

	var query string = `
	UPDATE blogs
	SET title = COALESCE($1, title),
		content = COALESCE($2, content),
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $3
	RETURNING id, title, content, created_at, updated_at
	`

	var blog models.Blog

	var err error = pool.QueryRow(ctx,query,title,content,id).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)

	if err != nil{
		return nil,err
	}

	return &blog,nil
}

func DeleteBlog(pool *pgxpool.Pool,id int)(*models.Blog,error){
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(),5 *time.Second)
	defer cancel()

	var query string = `
	DELETE FROM blogs
	WHERE ID = $1
	RETURNING id, title, content, created_at, updated_at
	`

	var blog models.Blog

	var err error = pool.QueryRow(ctx,query,id).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)

	if err != nil{
		return nil,err
	}

	return &blog,nil
}