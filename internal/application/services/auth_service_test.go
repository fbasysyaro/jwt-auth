package services_test

import (
	"context"
	"testing"

	"jwt-auth/internal/application/dto"
	appservices "jwt-auth/internal/application/services"
	"jwt-auth/internal/infrastructure/jwt"
)

func TestAuthService_Register(t *testing.T) {
	userRepo := newMockUserRepository()
	jwtManager := jwt.NewJWTManager()
	authService := appservices.NewAuthService(userRepo, jwtManager, nil, nil)

	t.Run("Valid registration", func(t *testing.T) {
		req := &dto.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		resp, err := authService.Register(context.Background(), req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp == nil || resp.AccessToken == "" {
			t.Error("expected valid response and token")
		}
	})

	t.Run("Duplicate email", func(t *testing.T) {
		req := &dto.RegisterRequest{
			Username: "testuser2",
			Email:    "test@example.com",
			Password: "password123",
		}
		_, err := authService.Register(context.Background(), req)
		if err == nil {
			t.Error("expected error for duplicate email")
		}
	})
}
