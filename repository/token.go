package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	redisClient *redis.Client
}

func NewTokenRepository() *TokenRepository {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})
	
	return &TokenRepository{
		redisClient: client,
	}
}

func (r *TokenRepository) GenerateToken(userID int, username, email, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(1 * time.Hour).Unix(), // 1 hour for access token
		"iat":      time.Now().Unix(),
		"type":     "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (r *TokenRepository) GenerateRefreshToken(userID int, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days for refresh token
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (r *TokenRepository) ValidateToken(tokenString, secret string) (jwt.MapClaims, error) {
	// Check if token is blacklisted in Redis
	ctx := context.Background()
	exists, err := r.redisClient.Exists(ctx, "blacklist:"+tokenString).Result()
	if err == nil && exists > 0 {
		return nil, fmt.Errorf("token is blacklisted")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (r *TokenRepository) BlacklistToken(tokenString string) {
	ctx := context.Background()
	// Blacklist for 24 hours
	r.redisClient.Set(ctx, "blacklist:"+tokenString, true, 24*time.Hour)
}