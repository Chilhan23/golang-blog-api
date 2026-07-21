package models

import "time"

type Like struct {
	ID        int       `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	BlogID    int       `json:"blog_id" db:"blog_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}