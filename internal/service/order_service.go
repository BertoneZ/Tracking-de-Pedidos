package service
import ("context"
	"tracking/internal/domain"
	"tracking/internal/dto"
	"tracking/internal/repository"
)
type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}
func (s *OrderService) CreateOrder(ctx context.Context, customerID string, req dto.CreateOrderRequest) (*domain.Order, error) {
	order := &domain.Order{
		CustomerID: customerID,
		Status:     "PENDING",
		OriginLat:  req.OriginLat,
		OriginLng:  req.OriginLng,
		DestLat:    req.DestLat,
		DestLng:    req.DestLng,
	}

	err := s.repo.Create(ctx, order)
	return order, err
}
func (s *OrderService) GetPendingOrders(ctx context.Context) ([]domain.Order, error) {
	return s.repo.GetPending(ctx)
}
func (s *OrderService) AcceptOrder(ctx context.Context, orderID string, driverID string) error {
	return s.repo.AcceptOrder(ctx, orderID, driverID)
}
func (s *OrderService) GetOrderById(ctx context.Context, id string) (*domain.Order, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *OrderService) CompleteOrder(ctx context.Context, orderID string, driverID string) error {
	return s.repo.CompleteOrder(ctx, orderID, driverID)
}
func (s *OrderService) GetUserHistory(ctx context.Context, userID string) ([]domain.Order, error) {
	return s.repo.GetHistory(ctx, userID)
}