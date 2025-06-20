package routes

import (
	"jwt-auth/internal/interfaces/http/handlers"
	"jwt-auth/internal/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	authHandler *handlers.AuthHandler,
	jwtMiddleware *middleware.JWTMiddleware,
) *gin.Engine {

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "JWT Auth Service is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")

	// Public routes (no authentication required)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected routes (authentication required)
	protected := v1.Group("/")
	protected.Use(jwtMiddleware.RequireAuth())
	{
		protected.GET("/profile", authHandler.Profile)
		protected.POST("/logout", authHandler.Logout)

		// Add more protected routes here
		protected.GET("/dashboard", func(c *gin.Context) {
			userID := c.GetInt("user_id")
			username := c.GetString("username")

			c.JSON(200, gin.H{
				"message":  "Welcome to dashboard",
				"user_id":  userID,
				"username": username,
			})
		})
	}

	return router
}
