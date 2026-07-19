package models

import "time"

type Blog struct {
	ID int `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"` 
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

