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
	orderRepo repository.OrderRepositoryInterface // Inyectamos órdenes para validar estado
}

func NewLocationService(repo repository.LocationRepositoryInterface, orderRepo repository.OrderRepositoryInterface) *LocationService {
	return &LocationService{
		repo:      repo,
		orderRepo: orderRepo,
	}
}

func (s *LocationService) UpdateLocation(ctx context.Context, driverID string, lat, lng float64) error {
	// 1. Validación de Coordenadas (Regla técnica)
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		return errors.New("coordenadas geográficas inválidas")
	}

	// 2. Validación de Negocio: ¿El driver tiene una orden activa?
	// Evitamos que drivers que no están trabajando saturen Redis con pings de GPS.
	active, err := s.orderRepo.HasActiveOrder(ctx, driverID)
	if err != nil {
		return err
	}
	if !active {
		slog.Warn("intento de update de ubicación de driver sin orden activa", "driver_id", driverID)
		return errors.New("no puedes reportar ubicación sin un pedido asignado")
	}

	// 3. Persistencia
	return s.repo.SaveDriverLocation(ctx, driverID, lat, lng)
}

func (s *LocationService) GetLocation(ctx context.Context, driverID string) (*redis.GeoLocation, error) {
	// Aquí se podría agregar lógica de "expiración" si la data es muy vieja
	return s.repo.GetDriverLocation(ctx, driverID)
}