package redis

import (
	"context"
	"fmt"
	"jwt-auth/internal/domain/services"
	"time"

	"github.com/redis/go-redis/v9"
)

type tokenBlacklistServiceImpl struct {
	redisClient *redis.Client
}

func NewTokenBlacklistService(redisClient *redis.Client) services.TokenBlacklistService {
	return &tokenBlacklistServiceImpl{
		redisClient: redisClient,
	}
}

func (s *tokenBlacklistServiceImpl) BlacklistToken(ctx context.Context, token string, expiresIn int64) error {
	key := fmt.Sprintf("blacklist:%s", token)
	return s.redisClient.Set(ctx, key, true, time.Duration(expiresIn)*time.Second).Err()
}

func (s *tokenBlacklistServiceImpl) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", token)
	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}