package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"errors"
)
type LocationRepositoryInterface interface {
	SaveDriverLocation(ctx context.Context, driverID string, lat, lng float64) error
	GetDriverLocation(ctx context.Context, driverID string) (*redis.GeoLocation, error)
	DeleteDriverLocation(ctx context.Context, driverID string) error
}
type LocationRepository struct {
	redis *redis.Client
}

func NewLocationRepository(r *redis.Client) *LocationRepository {
	return &LocationRepository{redis: r}
}

const DriversKey = "drivers_locations"

func (r *LocationRepository) SaveDriverLocation(ctx context.Context, driverID string, lat, lng float64) error {
	return r.redis.GeoAdd(ctx, DriversKey, &redis.GeoLocation{
		Name:      driverID,
		Latitude:  lat,
		Longitude: lng,
	}).Err()
}

func (r *LocationRepository) UpdateDriverLocation(ctx context.Context, driverID string, lat, lng float64) error {
	
	return r.redis.GeoAdd(ctx, "drivers_locations", &redis.GeoLocation{
		Name:      driverID,
		Latitude:  lat,
		Longitude: lng,
	}).Err()
}
func (r *LocationRepository) GetDriverLocation(ctx context.Context, driverID string) (*redis.GeoLocation, error) {
	pos, err := r.redis.GeoPos(ctx, "drivers_locations", driverID).Result()
	if err != nil || len(pos) == 0 || pos[0] == nil {
		return nil, errors.New("ubicación no encontrada para este repartidor")
	}

	return &redis.GeoLocation{
		Name:      driverID,
		Longitude: pos[0].Longitude,
		Latitude:  pos[0].Latitude,
	}, nil
}

func (r *LocationRepository) DeleteDriverLocation(ctx context.Context, driverID string) error {
	
	return r.redis.ZRem(ctx, DriversKey, driverID).Err()
}