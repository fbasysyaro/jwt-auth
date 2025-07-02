package main

import (
	"fmt"
	appservices "jwt-auth/internal/application/services"
	"jwt-auth/internal/infrastructure/database"
	emailinfra "jwt-auth/internal/infrastructure/email"
	"jwt-auth/internal/infrastructure/jwt"
	redisinfra "jwt-auth/internal/infrastructure/redis"
	"jwt-auth/internal/infrastructure/repositories"
	"jwt-auth/internal/interfaces/config"
	"jwt-auth/internal/interfaces/http/handlers"
	"jwt-auth/internal/interfaces/http/middleware"
	"jwt-auth/internal/interfaces/http/routes"
	"log"
	"os"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Redis
	redisClient := redisinfra.NewRedisClient(&redisinfra.RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	// Initialize token blacklist service
	tokenBlacklist := redisinfra.NewTokenBlacklistService(redisClient)

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
