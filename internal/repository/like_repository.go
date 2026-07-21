package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TogleLike(pool *pgxpool.Pool, UserID string, BlogID int) (bool, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	checkquery := `SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND blog_id = $2)`
	err := pool.QueryRow(ctx, checkquery, UserID, BlogID).Scan(&exists)
	if err != nil {
		return false, 0, err
	}

	var isLiked bool
	if exists {
		deletequery := `DELETE FROM likes WHERE user_id = $1 AND blog_id = $2`
		_, err := pool.Exec(ctx, deletequery, UserID, BlogID)
		if err != nil {
			return false, 0, err
		}
		isLiked = false
	} else {
		insertQuery := `INSERT INTO likes (user_id, blog_id) VALUES ($1, $2)`
		_, err = pool.Exec(ctx, insertQuery, UserID, BlogID)
		if err != nil {
			return false, 0, err
		}
		isLiked = true
	}

	var totalLikes int
	countQuery := `SELECT COUNT(*) FROM likes WHERE blog_id = $1`
	err = pool.QueryRow(ctx, countQuery, BlogID).Scan(&totalLikes)
	if err != nil {
		return isLiked, 0, err
	}
	return isLiked, totalLikes, nil
}