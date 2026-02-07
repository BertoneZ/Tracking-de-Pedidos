package dto
import "time"
type CreateOrderRequest struct {
	DestinationAddress string             `json:"destination_address" binding:"required"`
	Items              []OrderItemRequest `json:"items" binding:"required,gt=0"`
}
type OrderItemResponse struct {
    ProductName string  `json:"product_name"`
    Quantity    int     `json:"quantity"`
}
type OrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}
type OrderResponse struct {
	ID     string `json:"id"`
	DriverID string `json:"driver_id,omitempty"`
	CustomerID string `json:"customer_id"`
	CustomerName       string `json:"customer_name"`
	DestinationAddress string `json:"destination_address"`
	TotalPrice         float64 `json:"total_price"`
	Status string `json:"status"`
	Items  []OrderItemResponse `json:"items"`
	CreatedAt time.Time `json:"created_at"`
}
type UpdateLocationRequest struct {
	Lat float64 `json:"lat" binding:"required"`
	Lng float64 `json:"lng" binding:"required"`
}
