package models

import "time"

type User struct {
	ID string 			`json:"id" db:"id"`
	Username string 	`json:"username" db:"username"`
	Email string     	`json:"email" db:"email"`
	Password string 	`json:"-" db:"password"`
	Role string 		`json:"-" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"` 
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	
}

