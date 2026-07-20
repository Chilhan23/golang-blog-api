package models

import "time"

type Category struct {
	ID int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
	CreatedAt time.Time `json:"created_at" db:"created_at"` 
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

