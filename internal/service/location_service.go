package service
import (
	"tracking/internal/repository"
	"context"
	"github.com/redis/go-redis/v9"
)
type LocationService struct {
	repo *repository.LocationRepository
}
func NewLocationService(repo *repository.LocationRepository) *LocationService {
	return &LocationService{repo: repo}
}
func (s *LocationService) UpdateLocation(ctx context.Context, driverID string, lat, lng float64) error {
	// Aquí podrías validar, por ejemplo, que las coordenadas estén dentro de Rafaela
	return s.repo.UpdateDriverLocation(ctx, driverID, lat, lng)
}
func (s *LocationService) GetLocation(ctx context.Context, driverID string) (*redis.GeoLocation, error) {
	return s.repo.GetDriverLocation(ctx, driverID)
}