package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kaa-dan/clean-architecture-go/internal/infrastructure/security"
	"github.com/kaa-dan/clean-architecture-go/internal/interfaces/handlers"
	"github.com/kaa-dan/clean-architecture-go/pkg/response"
)

func SetupRoutes(
	router *gin.Engine,
	userHandler *handlers.UserHandler,
	authMiddleware *security.AuthMiddleware,
) {
	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(requestid.New())

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rate limiting middleware (simple implementation)
	router.Use(func(c *gin.Context) {
		// Add rate limiting logic here if needed
		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		response.Success(c, 200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
		})
	})

	// API routes
	api := router.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	{
		auth.POST("/signup", userHandler.SignUp)
		auth.POST("/signin", userHandler.SignIn)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(authMiddleware.RequireAuth())
	{
		// User profile routes
		protected.GET("/profile", userHandler.GetProfile)

		// User management routes
		users := protected.Group("/users")
		{
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Admin only routes
		admin := protected.Group("/admin")
		admin.Use(authMiddleware.RequireAdmin())
		{
			admin.GET("/users", userHandler.GetAllUsers)
		}
	}
}
