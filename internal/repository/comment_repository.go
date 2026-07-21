package repository

import (
	"context"
	"rest-api/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateComment(pool *pgxpool.Pool, BlogID int, UserID string, content string) (*models.Comment, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO comments (user_id, blog_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, blog_id, content, created_at, updated_at
	`

	var comment models.Comment
	var err error = pool.QueryRow(ctx, query, UserID, BlogID, content).Scan(
		&comment.ID,
		&comment.UserID,
		&comment.BlogID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func GetCommentsByBlogID(pool *pgxpool.Pool, blogID int) ([]models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT c.id, c.blog_id, c.user_id, u.username AS user_name, c.content, c.created_at, c.updated_at
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.blog_id = $1
		ORDER BY c.created_at ASC
	`
	rows, err := pool.Query(ctx, query, blogID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []models.Comment{}

	for rows.Next() {
		var comment models.Comment

		err := rows.Scan(
			&comment.ID,
			&comment.BlogID,
			&comment.UserID,
			&comment.UserName,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func DeleteComment(pool *pgxpool.Pool, commentID int, userID string, userRole string) (*models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		DELETE FROM comments
		WHERE id = $1 AND (user_id = $2 OR $3 = 'admin')
		RETURNING id, blog_id, user_id, content, created_at, updated_at
	`

	var comment models.Comment
	err := pool.QueryRow(ctx, query, commentID, userID, userRole).Scan(
		&comment.ID,
		&comment.BlogID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}