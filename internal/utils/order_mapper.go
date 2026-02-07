package utils

import (
	"tracking/internal/domain"
	"tracking/internal/dto"
)

func OrderDomainToResponseOrderDto(order domain.Order, customerName string, items []dto.OrderItemResponse) dto.OrderResponse {
	var itemsDto []dto.OrderItemResponse
	for _, item := range items {
		itemsDto = append(itemsDto, dto.OrderItemResponse{
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
		})
	}
	return dto.OrderResponse{
		ID:                 order.ID,
		DriverID:           order.DriverID,
		CustomerID:         order.CustomerID,
		CustomerName:       order.CustomerName,
		DestinationAddress: order.DestinationAddress,
		TotalPrice:         order.TotalPrice,
		Status:             order.Status,
		Items:              itemsDto,
		CreatedAt:          order.CreatedAt,
	}
}

func SliceOrderDomainToOrderResponseListDto(orders []domain.Order) []dto.OrderResponse {
	if orders == nil {
		return []dto.OrderResponse{}
	}

	res := make([]dto.OrderResponse, len(orders))
	for i, o := range orders {
		res[i] = OrderDomainToResponseOrderDto(o, "", nil)
	}
	return res
}


func ToOrderDomain(req dto.CreateOrderRequest, customerID string) *domain.Order {
	var items []domain.OrderItem
	for _, item := range req.Items {
		items = append(items, domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return &domain.Order{
		CustomerID:         customerID,
		DestinationAddress: req.DestinationAddress,
		Status:             "PENDING",
		Items:              items,
	}
}
