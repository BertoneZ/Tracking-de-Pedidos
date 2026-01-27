package domain

import "time"

type Order struct {
	ID          string    `json:"id"`
	CustomerID  string    `json:"customer_id"`
	DriverID    string   `json:"driver_id"` // Puntero porque puede ser nulo al inicio
	Status      string    `json:"status"`    // PENDING, ASSIGNED, etc.
	OriginLat   float64   `json:"origin_lat"`
	OriginLng   float64   `json:"origin_lng"`
	DestLat     float64   `json:"dest_lat"`
	DestLng     float64   `json:"dest_lng"`
	CreatedAt   time.Time `json:"created_at"`
}