package service

import (
	"context"
	
	"tracking/internal/dto"
	"tracking/internal/repository"
   
    "tracking/internal/utils"
    "log/slog"
)
const (
    //ubicacion fija del local, yo puse rafaela
    defaultOriginLat = -31.2503
    defaultOriginLng = -61.4867
)
type OrderServiceInterface interface {
    CreateOrder(ctx context.Context, req dto.CreateOrderRequest, customerID string) (string, error)
    GetPendingOrders(ctx context.Context) ([]dto.OrderResponse, error)
    AcceptOrder(ctx context.Context, orderID string, driverID string) error
    GetOrderById(ctx context.Context, id string) (dto.OrderResponse, error)
    CompleteOrder(ctx context.Context, orderID string, driverID string) error
    GetUserHistory(ctx context.Context, userID string) ([]dto.OrderResponse, error)
}
type OrderService struct {
	repo repository.OrderRepositoryInterface
    productRepo repository.ProductRepositoryInterface
}

func NewOrderService(repo repository.OrderRepositoryInterface, prodRepo repository.ProductRepositoryInterface) *OrderService {
	return &OrderService{
        repo: repo,
        productRepo: prodRepo,
    }
}
func (s *OrderService) CreateOrder(ctx context.Context, req dto.CreateOrderRequest, customerID string) (string, error) {
	order := utils.ToOrderDomain(req, customerID)

	lat, lng, err := s.GetCoordinates(order.DestinationAddress)
	if err != nil {
		slog.Error("error geocoding", "address", order.DestinationAddress, "error", err)
		return "", utils.ErrInvalidAddress
	}
	order.DestLat = lat
	order.DestLng = lng

	order.OriginLat = defaultOriginLat
	order.OriginLng = defaultOriginLng

	var totalPrice float64
	for i := range order.Items {
		product, err := s.productRepo.GetByID(ctx, order.Items[i].ProductID)
		if err != nil {
			return "", utils.ErrProductNotFound
		}
	
		order.Items[i].PriceAtTime = product.Price
		totalPrice += product.Price * float64(order.Items[i].Quantity)
	}

	order.TotalPrice = totalPrice

	return s.repo.CreateWithItems(ctx, order)
}
func (s *OrderService) GetPendingOrders(ctx context.Context) ([]dto.OrderResponse, error) {
    orders, err := s.repo.GetPending(ctx)
    if err != nil {
        slog.Error("error al obtener pedidos pendientes", "error", err)
        return nil, utils.ErrInternal 
    }
    return utils.SliceOrderDomainToOrderResponseListDto(orders), nil
}
func (s *OrderService) AcceptOrder(ctx context.Context, orderID string, driverID string) error {
	active, err := s.repo.HasActiveOrder(ctx, driverID)
    if err != nil {
        return err
    }
    if active {
        return utils.ErrDeliveryNotFinished
    }
    order, err := s.repo.GetOrderById(ctx, orderID)
    if err != nil {
        return err
    }
    if order.Status != "PENDING" {
        return utils.ErrOrderNotAvailable
    }
    return s.repo.AcceptOrder(ctx, orderID, driverID)
}
func (s *OrderService) GetOrderById(ctx context.Context, id string) (dto.OrderResponse, error) {
	order, err := s.repo.GetOrderById(ctx, id)
    if err != nil {
       slog.Error("pedido no encontrado", "id", id, "error", err)
        return dto.OrderResponse{}, utils.ErrOrderNotFound
    }
    return utils.OrderDomainToResponseOrderDto(order, "", nil), nil
}
func (s *OrderService) CompleteOrder(ctx context.Context, orderID string, driverID string) error {
        
    order, err := s.repo.GetOrderById(ctx, orderID)
    if err != nil {
        return utils.ErrOrderNotFound
    }

    if order.DriverID != driverID {
        slog.Warn("intento de completar orden ajena", "order_id", orderID, "driver_id", driverID)
        return utils.ErrUnauthorizedAction 
    }
    
    if order.Status != "ASSIGNED" {
        return utils.ErrInvalidState 
    }

    err = s.repo.CompleteOrder(ctx, orderID, driverID)
    if err != nil {
        slog.Error("error técnico al completar orden", "order_id", orderID, "error", err)
        return utils.ErrInternal
    }

    return nil
}
func (s *OrderService) GetUserHistory(ctx context.Context, userID string) ([]dto.OrderResponse, error) {
	orders, err := s.repo.GetHistory(ctx, userID)
   if err != nil {
        slog.Error("error al obtener historial", "user_id", userID, "error", err)  
        return nil, utils.ErrInternal 
    }
   
    if orders == nil {
        return []dto.OrderResponse{}, nil
    }
    return utils.SliceOrderDomainToOrderResponseListDto(orders), nil
}