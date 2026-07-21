package models

import "time"

type Blog struct {
	ID           int       `json:"id" db:"id"`
	Title        string    `json:"title" db:"title"`
	Content      string    `json:"content" db:"content"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	UserID       string    `json:"user_id" db:"user_id"`
	AuthorName   string    `json:"author_name,omitempty" db:"author_name"`
	CategoryID   *int      `json:"category_id" db:"category_id"`
	CategoryName *string   `json:"category_name,omitempty" db:"category_name"`
	TotalLikes   int       `json:"total_likes" db:"total_likes"`    
	IsLiked      bool      `json:"is_liked,omitempty" db:"is_liked"`
}
