package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Implements services.TokenBlacklistService
type TokenBlacklistService struct {
	redisClient *redis.Client
}

func NewTokenBlacklistService(redisClient *redis.Client) *TokenBlacklistService {
	return &TokenBlacklistService{
		redisClient: redisClient,
	}
}

func (s *TokenBlacklistService) BlacklistToken(ctx context.Context, token string, expiresIn int64) error {
	key := fmt.Sprintf("blacklist:%s", token)
	return s.redisClient.Set(ctx, key, true, time.Duration(expiresIn)*time.Second).Err()
}

func (s *TokenBlacklistService) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", token)
	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
