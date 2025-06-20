package main

import (
	"fmt"
	"jwt-auth/internal/application/services"
	"jwt-auth/internal/infrastructure/database"
	"jwt-auth/internal/infrastructure/jwt"
	"jwt-auth/internal/infrastructure/repositories"
	"jwt-auth/internal/interfaces/config"
	"jwt-auth/internal/interfaces/http/handlers"
	"jwt-auth/internal/interfaces/http/middleware"
	"jwt-auth/internal/interfaces/http/routes"
	"log"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

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
	jwtManager := jwt.NewJWTManager(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenExpiry,
		cfg.JWT.RefreshTokenExpiry,
	)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, jwtManager)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize middleware
	jwtMiddleware := middleware.NewJWTMiddleware(authService)

	// Setup routes
	router := routes.SetupRoutes(authHandler, jwtMiddleware)

	// Start server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
