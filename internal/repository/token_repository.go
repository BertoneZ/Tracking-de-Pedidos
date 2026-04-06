package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepositoryInterface interface {
	StoreRefreshToken(ctx context.Context, userID, refreshToken string, ttl time.Duration) error
	IsValidRefreshToken(ctx context.Context, userID, refreshToken string) (bool, error)
	DeleteRefreshToken(ctx context.Context, userID string) error
}

type TokenRepository struct {
	rdb *redis.Client
}

func NewTokenRepository(rdb *redis.Client) *TokenRepository {
	return &TokenRepository{rdb: rdb}
}

func (r *TokenRepository) StoreRefreshToken(ctx context.Context, userID, refreshToken string, ttl time.Duration) error {
	key := refreshTokenKey(userID)
	return r.rdb.Set(ctx, key, refreshToken, ttl).Err()
}

func (r *TokenRepository) IsValidRefreshToken(ctx context.Context, userID, refreshToken string) (bool, error) {
	key := refreshTokenKey(userID)
	storedToken, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return storedToken == refreshToken, nil
}

func (r *TokenRepository) DeleteRefreshToken(ctx context.Context, userID string) error {
	key := refreshTokenKey(userID)
	if err := r.rdb.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func refreshTokenKey(userID string) string {
	return fmt.Sprintf("auth:refresh:%s", userID)
}
