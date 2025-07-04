// ...existing code...
package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"jwt-auth/internal/application/dto"
	"jwt-auth/internal/domain/entities"
	"jwt-auth/internal/domain/repositories"
	"jwt-auth/internal/domain/services"

	"golang.org/x/crypto/bcrypt"
)

type authServiceImpl struct {
	userRepo       repositories.UserRepository
	jwtManager     services.JWTManager
	emailService   services.EmailService
	tokenBlacklist services.TokenBlacklistService
}

func NewAuthService(userRepo repositories.UserRepository, jwtManager services.JWTManager, emailService services.EmailService, tokenBlacklist services.TokenBlacklistService) services.AuthService {
	return &authServiceImpl{
		userRepo:       userRepo,
		jwtManager:     jwtManager,
		emailService:   emailService,
		tokenBlacklist: tokenBlacklist,
	}
}

func (s *authServiceImpl) Logout(ctx context.Context, token string) error {
	// Blacklist the token using the injected tokenBlacklist service
	// For demo, use 1 hour expiration
	if s.tokenBlacklist != nil {
		err := s.tokenBlacklist.BlacklistToken(ctx, token, int64(time.Hour.Seconds()))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *authServiceImpl) InitiatePasswordReset(ctx context.Context, email string) error {
	// Check if user exists
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Generate reset token
	resetToken, err := s.jwtManager.GenerateToken(fmt.Sprintf("%d", user.ID), map[string]interface{}{
		"email": user.Email,
		"type":  "password_reset",
	})
	if err != nil {
		return err
	}

	// Send reset email
	subject := "Password Reset Request"
	body := fmt.Sprintf("Click here to reset your password: http://localhost:8080/api/v1/auth/reset-password?token=%s", resetToken)
	return s.emailService.SendEmail(ctx, email, subject, body)
}

func (s *authServiceImpl) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate reset token
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid reset token")
	}

	// Check token type
	if tokenType, ok := claims["type"].(string); !ok || tokenType != "password_reset" {
		return fmt.Errorf("invalid token type")
	}

	// Get user
	userID, _ := claims["user_id"].(string)
	userIntID := 0
	fmt.Sscanf(userID, "%d", &userIntID)
	user, err := s.userRepo.GetByID(ctx, userIntID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	user.Password = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}

func (s *authServiceImpl) VerifyEmail(ctx context.Context, token string) error {
	// Validate verification token
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid verification token")
	}

	// Check token type
	if tokenType, ok := claims["type"].(string); !ok || tokenType != "email_verification" {
		return fmt.Errorf("invalid token type")
	}

	// Get user
	userID, _ := claims["user_id"].(string)
	userIntID := 0
	fmt.Sscanf(userID, "%d", &userIntID)
	user, err := s.userRepo.GetByID(ctx, userIntID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Mark email as verified
	user.EmailVerified = true
	return s.userRepo.Update(ctx, user)
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

	// Send email verification
	verificationToken, _ := s.jwtManager.GenerateToken(fmt.Sprintf("%d", user.ID), map[string]interface{}{
		"email": user.Email,
		"type":  "email_verification",
	})
	subject := "Email Verification"
	body := fmt.Sprintf("Click here to verify your email: http://localhost:8080/api/v1/auth/verify-email/%s", verificationToken)
	s.emailService.SendEmail(ctx, user.Email, subject, body)

	// Generate tokens using domain interface
	claims := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	}
	userID := fmt.Sprintf("%d", user.ID)
	accessToken, err := s.jwtManager.GenerateToken(userID, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err := s.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
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

	// Generate tokens using domain interface
	claims := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	}
	userID := fmt.Sprintf("%d", user.ID)
	accessToken, err := s.jwtManager.GenerateToken(userID, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err := s.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
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
	// Validate refresh token
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token claims")
	}
	// Generate new tokens
	accessToken, err := s.jwtManager.GenerateToken(userID, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	// Get user details
	userIntID := 0
	fmt.Sscanf(userID, "%d", &userIntID)
	user, err := s.userRepo.GetByID(ctx, userIntID)
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
	// Check if token is blacklisted
	if s.tokenBlacklist != nil {
		blacklisted, err := s.tokenBlacklist.IsTokenBlacklisted(ctx, token)
		if err != nil {
			return nil, err
		}
		if blacklisted {
			return nil, fmt.Errorf("token is blacklisted")
		}
	}
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	userID, _ := claims["user_id"].(string)
	username, _ := claims["username"].(string)
	email, _ := claims["email"].(string)
	userIntID := 0
	fmt.Sscanf(userID, "%d", &userIntID)
	return &dto.UserClaims{
		UserID:   userIntID,
		Username: username,
		Email:    email,
	}, nil
}
