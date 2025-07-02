package services_test

import (
	"context"
	"testing"

	"jwt-auth/internal/application/dto"
	appservices "jwt-auth/internal/application/services"
	"jwt-auth/internal/infrastructure/jwt"
	redisService "jwt-auth/internal/infrastructure/redis"

	"github.com/redis/go-redis/v9"
)

func TestAuthService_Integration(t *testing.T) {
	// Setup dependencies
	userRepo := newMockUserRepository()
	jwtManager := jwt.NewJWTManager()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     ":6379",                   // Use a real or mock Redis instance
		Password: "RedisSecurePassword123!", // Use the password from your .env/docker-compose
	})
	emailSvc := newMockEmailService()
	tokenBlacklist := redisService.NewTokenBlacklistService(redisClient)

	authService := appservices.NewAuthService(userRepo, jwtManager, emailSvc, tokenBlacklist)

	t.Run("Full auth flow", func(t *testing.T) {
		ctx := context.Background()

		// 1. Register
		registerReq := &dto.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		regResp, err := authService.Register(ctx, registerReq)
		if err != nil {
			t.Fatalf("Register failed: %v", err)
		}
		if regResp.AccessToken == "" || regResp.RefreshToken == "" {
			t.Error("Register should return valid tokens")
		}

		// 2. Login
		loginReq := &dto.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		loginResp, err := authService.Login(ctx, loginReq)
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		if loginResp.AccessToken == "" || loginResp.RefreshToken == "" {
			t.Error("Login should return valid tokens")
		}

		// 3. Validate Token
		userClaims, err := authService.ValidateToken(ctx, loginResp.AccessToken)
		if err != nil {
			t.Fatalf("Token validation failed: %v", err)
		}
		if userClaims.Email != registerReq.Email {
			t.Errorf("Token validation returned wrong user. got: %s, want: %s", userClaims.Email, registerReq.Email)
		}

		// 4. Refresh Token
		refreshResp, err := authService.RefreshToken(ctx, loginResp.RefreshToken)
		if err != nil {
			t.Fatalf("Token refresh failed: %v", err)
		}
		if refreshResp.AccessToken == "" || refreshResp.RefreshToken == "" {
			t.Error("Refresh should return new valid tokens")
		}

		// 5. Logout
		if err := authService.Logout(ctx, loginResp.AccessToken); err != nil {
			t.Fatalf("Logout failed: %v", err)
		}

		// Verify token is blacklisted
		if _, err := authService.ValidateToken(ctx, loginResp.AccessToken); err == nil {
			t.Error("Should not be able to use token after logout")
		}
	})
}
