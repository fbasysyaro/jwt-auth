package services

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBlacklistService struct {
	redisClient *redis.Client
}

func NewTokenBlacklistService(redisClient *redis.Client) *TokenBlacklistService {
	return &TokenBlacklistService{
		redisClient: redisClient,
	}
}

func (s *TokenBlacklistService) BlacklistToken(ctx context.Context, token string, expiresIn time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", token)
	return s.redisClient.Set(ctx, key, true, expiresIn).Err()
}

func (s *TokenBlacklistService) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", token)
	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
