package dto

type CreateOrderRequest struct {
	OriginLat   float64 `json:"origin_lat" binding:"required"`
	OriginLng   float64 `json:"origin_lng" binding:"required"`
	DestLat     float64 `json:"dest_lat" binding:"required"`
	DestLng     float64 `json:"dest_lng" binding:"required"`
}

type OrderResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
type UpdateLocationRequest struct {
	Lat float64 `json:"lat" binding:"required"`
	Lng float64 `json:"lng" binding:"required"`
}