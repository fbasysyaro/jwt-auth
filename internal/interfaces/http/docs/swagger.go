// Package docs JWT Authentication Service API Documentation
//
// This is a JWT authentication service that provides user registration, login, and token management.
//
//	Schemes: http, https
//	Host: localhost:8080
//	BasePath: /api/v1
//	Version: 1.0.0
//	Contact: API Support<support@example.com>
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- Bearer: []
//
// swagger:meta
package docs

import (
	"jwt-auth/internal/application/dto"
)

// swagger:route POST /auth/register auth register
// Register a new user.
// responses:
//   200: authResponse
//   400: errorResponse

// swagger:route POST /auth/login auth login
// Login with email and password.
// responses:
//   200: authResponse
//   401: errorResponse

// swagger:route POST /auth/refresh auth refreshToken
// Refresh access token using refresh token.
// responses:
//   200: authResponse
//   401: errorResponse

// swagger:route POST /auth/reset-password auth resetPassword
// Request password reset email.
// responses:
//   200: successResponse
//   400: errorResponse

// swagger:route POST /auth/verify-email auth verifyEmail
// Verify email address.
// responses:
//   200: successResponse
//   400: errorResponse

// swagger:route GET /profile profile getProfile
// Get user profile information.
// Security:
//   - Bearer: []
// responses:
//   200: userResponse
//   401: errorResponse

// swagger:parameters register
type registerParams struct {
	// User registration data
	// in:body
	Body dto.RegisterRequest
}

// swagger:parameters login
type loginParams struct {
	// User login credentials
	// in:body
	Body dto.LoginRequest
}

// swagger:response authResponse
type authResponseWrapper struct {
	// in:body
	Body dto.AuthResponse
}

// swagger:response errorResponse
type errorResponseWrapper struct {
	// in:body
	Body dto.ErrorResponse
}

// swagger:response successResponse
type successResponseWrapper struct {
	// in:body
	Body dto.SuccessResponse
}

// swagger:response userResponse
type userResponseWrapper struct {
	// in:body
	Body struct {
		User struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"user"`
	}
}
