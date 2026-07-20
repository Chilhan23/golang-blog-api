package models

import "time"

type Blog struct {
	ID           int       `json:"id" db:"id"`
	Title        string    `json:"title" db:"title"`
	Content      string    `json:"content" db:"content"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	UserID       string    `json:"user_id" db:"user_id"`
	CategoryID   *int      `json:"category_id" db:"category_id"`
	CategoryName *string   `json:"category_name,omitempty" db:"category_name"`
}
