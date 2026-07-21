package models

import "time"

type Comment struct{
	ID        int       `json:"id" db:"id"`
	BlogID int `json:"blog_id" db:"blog_id"`
	UserID string `json:"user_id" db:"user_id"`
	Content string `json:"content" db:"content"`
	UserName  string    `json:"user_name,omitempty" db:"user_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"` 
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}