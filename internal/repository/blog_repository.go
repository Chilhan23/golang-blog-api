package repository

import (
	"context"
	"rest-api/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateBlog(pool *pgxpool.Pool, title string, content string,userID string)(*models.Blog, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	
	ctx, cancel = context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO blogs (title, content,user_id ,created_at, updated_at)
		VALUES ($1, $2, $3 , NOW(), NOW())
		RETURNING id, title, content, user_id , created_at, updated_at
	`
	var blog models.Blog
	var err error = pool.QueryRow(ctx,query,title,content,userID).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.UserID,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func GetAllBlogs(pool *pgxpool.Pool,userID string)([]models.Blog,error){
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	var query = `
		SELECT id,title,content,created_at,updated_at,user_id
		FROM blogs
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	var rows, err = pool.Query(ctx,query,userID)
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
			&blog.UserID,
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

func GetBlogByID(pool *pgxpool.Pool, id int,userID string)(*models.Blog,error){
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(),5 *time.Second)
	defer cancel()

	var query = `
		SELECT id,title,content,created_at,updated_at,user_id
		FROM blogs
		WHERE ID = $1 AND user_id = $2
	`
	var blog models.Blog

	err := pool.QueryRow(ctx,query,id,userID).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.CreatedAt,
		&blog.UpdatedAt,
		&blog.UserID,
	)

	if err != nil{
		return nil, err
	}

	return &blog,nil

}

func UpdateBlog(pool *pgxpool.Pool,id int,title string,content string,userID string)(*models.Blog,error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(),5 *time.Second)
	defer cancel()

	var query string = `
	UPDATE blogs
	SET title = COALESCE($1, title),
		content = COALESCE($2, content),
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $3 AND user_id = $4
	RETURNING id, title, content, created_at, updated_at,user_id
	`

	var blog models.Blog

	var err error = pool.QueryRow(ctx,query,title,content,id,userID).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.CreatedAt,
		&blog.UpdatedAt,
		&blog.UserID,
	)

	if err != nil{
		return nil,err
	}

	return &blog,nil
}

func DeleteBlog(pool *pgxpool.Pool,id int,userID string)(*models.Blog,error){
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(),5 *time.Second)
	defer cancel()

	var query string = `
	DELETE FROM blogs
	WHERE ID = $1 AND user_id = $2
	RETURNING id, title, content, created_at, updated_at, user_id
	`

	var blog models.Blog

	var err error = pool.QueryRow(ctx,query,id,userID).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.CreatedAt,
		&blog.UpdatedAt,
		&blog.UserID,
	)

	if err != nil{
		return nil,err
	}

	return &blog,nil
}