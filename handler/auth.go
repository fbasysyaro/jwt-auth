package handler

import (
	"jwt-auth/model"
	"jwt-auth/repository"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.TokenRepository
	jwtSecret string
}

func NewAuthHandler(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtSecret: jwtSecret,
	}
}

// Register user
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	if existingUser, _ := h.userRepo.GetByEmail(req.Email); existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Create user
	user, err := h.userRepo.Create(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate tokens
	accessToken, err := h.tokenRepo.GenerateToken(user.ID, user.Username, user.Email, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := h.tokenRepo.GenerateRefreshToken(user.ID, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusCreated, model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User:         user,
	})
}

// Login user - requires username/password in body + access token in Authorization header
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get access token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token required in Authorization header"})
		return
	}
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	// Validate access token
	claims, err := h.tokenRepo.ValidateToken(accessToken, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
		return
	}

	// Get user by username from token claims
	tokenUsername, _ := claims["username"].(string)
	if tokenUsername != req.Username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username mismatch"})
		return
	}

	// Get user from database
	tokenEmail, _ := claims["email"].(string)
	user, err := h.userRepo.GetByEmail(tokenEmail)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Validate password
	if !h.userRepo.ValidatePassword(user, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Login successful
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// Verify token
func (h *AuthHandler) Verify(c *gin.Context) {
	var req model.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Try to get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token required"})
			return
		}
		req.Token = strings.TrimPrefix(authHeader, "Bearer ")
	}

	// Validate token
	claims, err := h.tokenRepo.ValidateToken(req.Token, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Extract user info from claims
	userID, _ := claims["user_id"].(float64)
	username, _ := claims["username"].(string)
	email, _ := claims["email"].(string)

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"user": gin.H{
			"id":       int(userID),
			"username": username,
			"email":    email,
		},
	})
}

// Refresh token
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req model.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate refresh token
	claims, err := h.tokenRepo.ValidateToken(req.RefreshToken, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Check token type
	if tokenType, ok := claims["type"].(string); !ok || tokenType != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
		return
	}

	// Get user info
	userID, _ := claims["user_id"].(float64)
	user, err := h.userRepo.GetByID(int(userID))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Generate new tokens
	accessToken, err := h.tokenRepo.GenerateToken(user.ID, user.Username, user.Email, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	newRefreshToken, err := h.tokenRepo.GenerateRefreshToken(user.ID, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Blacklist old refresh token
	h.tokenRepo.BlacklistToken(req.RefreshToken)

	c.JSON(http.StatusOK, model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User:         user,
	})
}

// Logout user
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token required"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Blacklist token
	h.tokenRepo.BlacklistToken(token)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}