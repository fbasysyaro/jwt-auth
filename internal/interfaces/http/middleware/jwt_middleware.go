package middleware

import (
	"jwt-auth/internal/application/dto"
	"jwt-auth/internal/domain/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type JWTMiddleware struct {
	authService services.AuthService
}

func NewJWTMiddleware(authService services.AuthService) *JWTMiddleware {
	return &JWTMiddleware{
		authService: authService,
	}
}

func (m *JWTMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "unauthorized",
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check Bearer token format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "unauthorized",
				Message: "Authorization header must be in Bearer format",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		userClaims, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid token",
			})
			c.Abort()
			return
		}

		// Set user claims in context
		c.Set("user_id", userClaims.UserID)
		c.Set("username", userClaims.Username)
		c.Set("email", userClaims.Email)
		c.Set("user_claims", userClaims)

		c.Next()
	}
}

// Optional middleware for routes that may or may not require authentication
func (m *JWTMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")

			if userClaims, err := m.authService.ValidateToken(c.Request.Context(), token); err == nil {
				c.Set("user_id", userClaims.UserID)
				c.Set("username", userClaims.Username)
				c.Set("email", userClaims.Email)
				c.Set("user_claims", userClaims)
			}
		}

		c.Next()
	}
}
