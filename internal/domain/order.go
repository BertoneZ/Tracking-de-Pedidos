package domain

import "time"

type OrderItem struct {
	ProductID   string  `json:"product_id"`
    ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	PriceAtTime float64 `json:"price_at_time"`
}

type Order struct {
    ID                 string    `json:"id"`
    CustomerID         string    `json:"customer_id"`
    CustomerName       string    `json:"customer_name"` 
    DriverID           string    `json:"driver_id"`
    DriverName         string    `json:"driver_name"`   
    Status             string    `json:"status"`
    OriginLat          float64   `json:"origin_lat"`
    OriginLng          float64   `json:"origin_lng"`
    DestLat            float64   `json:"dest_lat"`
    DestLng            float64   `json:"dest_lng"`
    DestinationAddress string    `json:"destination_address"`
    TotalPrice         float64   `json:"total_price"`
    CreatedAt          time.Time `json:"created_at"`
    Items              []OrderItem `json:"items"`
}