package handlers

import (
	"jwt-auth/internal/application/dto"
	"jwt-auth/internal/domain/services"
	"jwt-auth/internal/interfaces/http/middleware"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration request"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	middleware.ValidateRequest(&req)(c)
	if c.IsAborted() {
		return
	}

	response, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "registration_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.AuthResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	middleware.ValidateRequest(&req)(c)
	if c.IsAborted() {
		return
	}

	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "login_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{refresh_token=string} true "Refresh token request"
// @Success 200 {object} dto.AuthResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	response, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "refresh_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Profile godoc
// @Summary Get user profile
// @Description Get current user profile information
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /profile [get]
func (h *AuthHandler) Profile(c *gin.Context) {
	userClaims, exists := c.Get("user_claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Profile retrieved successfully",
		Data:    userClaims,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout current user and blacklist token
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get the token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_token",
			Message: "Invalid token format",
		})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Process logout
	err := h.authService.Logout(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "logout_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Successfully logged out",
	})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	middleware.ValidateRequest(&req)(c)
	if c.IsAborted() {
		return
	}

	if err := h.authService.InitiatePasswordReset(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "password_reset_failed",
			Message: "Failed to initiate password reset",
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Password reset email sent",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}
	middleware.ValidateRequest(&req)(c)
	if c.IsAborted() {
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "reset_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Password reset successful",
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_token",
			Message: "Verification token is required",
		})
		return
	}

	if err := h.authService.VerifyEmail(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "verification_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Email verified successfully",
	})
}
