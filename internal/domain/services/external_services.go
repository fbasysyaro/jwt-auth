package services

import "context"

// JWTManager defines the interface for JWT operations
type JWTManager interface {
	GenerateToken(userID string, claims map[string]interface{}) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(token string) (map[string]interface{}, error)
}

// EmailService defines the interface for sending emails
type EmailService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}
