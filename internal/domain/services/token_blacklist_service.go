package services

import "context"

// TokenBlacklistService defines the interface for token blacklisting
type TokenBlacklistService interface {
	BlacklistToken(ctx context.Context, token string, expiration int64) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}