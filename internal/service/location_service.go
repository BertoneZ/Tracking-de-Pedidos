package service
import (
	"tracking/internal/repository"
	"context"
	"github.com/redis/go-redis/v9"
	"errors"
	"log/slog"
)
type LocationServiceInterface interface {
	UpdateLocation(ctx context.Context, driverID string, lat, lng float64) error
	GetLocation(ctx context.Context, driverID string) (*redis.GeoLocation, error)
}

type LocationService struct {
	repo      repository.LocationRepositoryInterface
	orderRepo repository.OrderRepositoryInterface 
}

func NewLocationService(repo repository.LocationRepositoryInterface, orderRepo repository.OrderRepositoryInterface) *LocationService {
	return &LocationService{
		repo:      repo,
		orderRepo: orderRepo,
	}
}

func (s *LocationService) UpdateLocation(ctx context.Context, driverID string, lat, lng float64) error {
	
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		return errors.New("coordenadas geográficas inválidas")
	}

	
	active, err := s.orderRepo.HasActiveOrder(ctx, driverID)
	if err != nil {
		return err
	}
	if !active {
		slog.Warn("intento de update de ubicación de driver sin orden activa", "driver_id", driverID)
		return errors.New("no puedes reportar ubicación sin un pedido asignado")
	}

	
	return s.repo.SaveDriverLocation(ctx, driverID, lat, lng)
}

func (s *LocationService) GetLocation(ctx context.Context, driverID string) (*redis.GeoLocation, error) {
	
	return s.repo.GetDriverLocation(ctx, driverID)
}