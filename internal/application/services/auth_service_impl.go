package services

import (
	"context"
	"fmt"
	"jwt-auth/internal/application/dto"
	"jwt-auth/internal/domain/entities"
	"jwt-auth/internal/domain/repositories"
	"jwt-auth/internal/domain/services"
	"jwt-auth/internal/infrastructure/jwt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authServiceImpl struct {
	userRepo   repositories.UserRepository
	jwtManager *jwt.JWTManager
}

func NewAuthService(userRepo repositories.UserRepository, jwtManager *jwt.JWTManager) services.AuthService {
	return &authServiceImpl{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *authServiceImpl) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &entities.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	userClaims := &dto.UserClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokens(userClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Hour.Seconds()), // 1 hour
		User:         user,
	}, nil
}

func (s *authServiceImpl) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate tokens
	userClaims := &dto.UserClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokens(userClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Hour.Seconds()), // 1 hour
		User:         user,
	}, nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	// Validate and refresh token
	accessToken, newRefreshToken, err := s.jwtManager.RefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Validate the new access token to get user claims
	userClaims, err := s.jwtManager.ValidateToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate new token: %w", err)
	}

	// Get user details
	user, err := s.userRepo.GetByID(ctx, userClaims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Hour.Seconds()), // 1 hour
		User:         user,
	}, nil
}

func (s *authServiceImpl) ValidateToken(ctx context.Context, token string) (*dto.UserClaims, error) {
	// Remove Bearer prefix if present
	token = strings.TrimPrefix(token, "Bearer ")

	return s.jwtManager.ValidateToken(token)
}
