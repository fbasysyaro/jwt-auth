package services

import (
	"context"
	"jwt-auth/internal/application/dto"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*dto.UserClaims, error)
}
