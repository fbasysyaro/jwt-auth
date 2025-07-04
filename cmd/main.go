// Package main JWT Authentication Service
//
// JWT Authentication Service API
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
//	SecurityDefinitions:
//	  Bearer:
//	    type: apiKey
//	    name: Authorization
//	    in: header
//	    description: Type "Bearer" followed by a space and JWT token.
//
// swagger:meta
package main

import (
	"fmt"
	appservices "jwt-auth/internal/application/services"
	"jwt-auth/internal/infrastructure/database"
	emailinfra "jwt-auth/internal/infrastructure/email"
	"jwt-auth/internal/infrastructure/jwt"
	"jwt-auth/internal/infrastructure/redis"
	"jwt-auth/internal/infrastructure/repositories"
	"jwt-auth/internal/interfaces/config"
	"jwt-auth/internal/interfaces/http/handlers"
	"jwt-auth/internal/interfaces/http/middleware"
	"jwt-auth/internal/interfaces/http/routes"
	"log"
	"os"

	_ "jwt-auth/docs"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Redis
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	redisClient := redis.NewRedisClient(&redis.RedisConfig{
		Host:     redisHost,
		Port:     redisPort,
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	// Initialize token blacklist service
	tokenBlacklist := redis.NewTokenBlacklistService(redisClient)

	// Initialize email service (for password reset)
	emailService := emailinfra.NewEmailService()

	// Initialize database
	db, err := database.NewPostgresDB(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	// Initialize JWT manager
	jwtManager := jwt.NewJWTManager()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	// Ensure infrastructure implementations are passed as domain interfaces

	authService := appservices.NewAuthService(
		userRepo,
		jwtManager,
		emailService,
		tokenBlacklist,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize middleware
	jwtMiddleware := middleware.NewJWTMiddleware(authService)

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(redisClient, 5, 60) // 100 requests per 60 seconds

	// Setup routes
	router := routes.SetupRoutes(authHandler, jwtMiddleware, rateLimiter)

	// Start server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
