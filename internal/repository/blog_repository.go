package repository

import (
	"context"
	"rest-api/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateBlog(pool *pgxpool.Pool, title string, content string, userID string, categoryID *int) (*models.Blog, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO blogs (title, content, user_id, category_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, title, content, user_id, category_id, created_at, updated_at
	`
	var blog models.Blog
	var err error = pool.QueryRow(ctx, query, title, content, userID, categoryID).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.UserID,
		&blog.CategoryID,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func GetAllBlogs(pool *pgxpool.Pool) ([]models.Blog, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query = `
		SELECT b.id, b.title, b.content, b.created_at, b.updated_at, b.user_id, u.username AS author_name, b.category_id, c.name AS category_name, COUNT(l.id) AS total_likes
		FROM blogs b
		LEFT JOIN users u ON b.user_id = u.id
		LEFT JOIN categories c ON b.category_id = c.id
		LEFT JOIN likes l ON b.id = l.blog_id
		GROUP BY b.id, u.username, c.name
		ORDER BY b.created_at DESC
	`
	var rows, err = pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []models.Blog = []models.Blog{}

	for rows.Next() {
		var blog models.Blog

		err = rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Content,
			&blog.CreatedAt,
			&blog.UpdatedAt,
			&blog.UserID,
			&blog.AuthorName,
			&blog.CategoryID,
			&blog.CategoryName,
			&blog.TotalLikes,
		)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}

func GetBlogsByUserID(pool *pgxpool.Pool, userID string) ([]models.Blog, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query = `
		SELECT b.id, b.title, b.content, b.created_at, b.updated_at, b.user_id, u.username AS author_name, b.category_id, c.name AS category_name, COUNT(l.id) AS total_likes
		FROM blogs b
		LEFT JOIN users u ON b.user_id = u.id
		LEFT JOIN categories c ON b.category_id = c.id
		LEFT JOIN likes l ON b.id = l.blog_id
		WHERE b.user_id = $1
		GROUP BY b.id, u.username, c.name
		ORDER BY b.created_at DESC
	`
	var rows, err = pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []models.Blog = []models.Blog{}

	for rows.Next() {
		var blog models.Blog

		err = rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Content,
			&blog.CreatedAt,
			&blog.UpdatedAt,
			&blog.UserID,
			&blog.AuthorName,
			&blog.CategoryID,
			&blog.CategoryName,
			&blog.TotalLikes,
		)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}

func GetBlogByID(pool *pgxpool.Pool, id int) (*models.Blog, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query = `
		SELECT b.id, b.title, b.content, b.created_at, b.updated_at, b.user_id, u.username AS author_name, b.category_id, c.name AS category_name, COUNT(l.id) AS total_likes
		FROM blogs b
		LEFT JOIN users u ON b.user_id = u.id
		LEFT JOIN categories c ON b.category_id = c.id
		LEFT JOIN likes l ON b.id = l.blog_id
		WHERE b.id = $1
		GROUP BY b.id, u.username, c.name
	`
	var blog models.Blog

	err := pool.QueryRow(ctx, query, id).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.CreatedAt,
		&blog.UpdatedAt,
		&blog.UserID,
		&blog.AuthorName,
		&blog.CategoryID,
		&blog.CategoryName,
		&blog.TotalLikes,
	)

	if err != nil {
		return nil, err
	}

	return &blog, nil
}

func UpdateBlog(pool *pgxpool.Pool, id int, title string, content string, userID string, categoryID *int) (*models.Blog, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		UPDATE blogs
		SET title = COALESCE($1, title),
			content = COALESCE($2, content),
			category_id = COALESCE($3, category_id),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND user_id = $5
		RETURNING id, title, content, user_id, category_id, created_at, updated_at
	`

	var blog models.Blog

	var err error = pool.QueryRow(ctx, query, title, content, categoryID, id, userID).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.UserID,
		&blog.CategoryID,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &blog, nil
}

func DeleteBlog(pool *pgxpool.Pool, id int, userID string) (*models.Blog, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM blogs
		WHERE id = $1 AND user_id = $2
		RETURNING id, title, content, user_id, category_id, created_at, updated_at
	`

	var blog models.Blog

	var err error = pool.QueryRow(ctx, query, id, userID).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.UserID,
		&blog.CategoryID,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &blog, nil
}