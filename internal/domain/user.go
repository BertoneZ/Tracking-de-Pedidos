package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	FullName	 string    `json:"full_name"`
	PasswordHash string    `json:"-"` // El guion hace que nunca se envíe en los JSON
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}