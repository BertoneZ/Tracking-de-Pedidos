package domain

import "time"

type Driver struct {
	UserID       string    `json:"user_id"`
	IsActive     bool      `json:"is_active"`
	LastLat      float64   `json:"lat"`
	LastLng      float64   `json:"lng"`
	UpdatedAt    time.Time `json:"updated_at"`
}