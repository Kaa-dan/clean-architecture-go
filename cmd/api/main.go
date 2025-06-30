package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kaa-dan/clean-architecture-go/internal/config"
	"github.com/kaa-dan/clean-architecture-go/internal/infrastructure/database"
	"github.com/kaa-dan/clean-architecture-go/internal/infrastructure/repositories"
	"github.com/kaa-dan/clean-architecture-go/internal/infrastructure/security"
	"github.com/kaa-dan/clean-architecture-go/internal/interfaces/handlers"
	"github.com/kaa-dan/clean-architecture-go/internal/interfaces/routes"
	"github.com/kaa-dan/clean-architecture-go/internal/usecases"
	"github.com/kaa-dan/clean-architecture-go/pkg/logger"
)

func main() {
	// Load configuration

	cfg := config.Load()

	// Initialize logger
	logger.Init(cfg.LogLevel)

	//Connect to MongoDb

	db, err := database.NewMongoDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer db.Disconnect(context.Background())

	// Initialize repositories

	userRepo := repositories.NewUserRepository(db, cfg.DatabaseName)

	// Initialize use cases
	jwtManager := security.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiryHours)
	passwordManager := security.NewPasswordManger()
	userUseCases := usecases.NewUserUseCase(userRepo, jwtManager, passwordManager)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUseCases)

	// Initialize middleware
	authMiddleware := security.NewAuthMiddleware(jwtManager)

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.New()

	// Setup routes
	routes.SetupRoutes(router, userHandler, authMiddleware)

	// Create server
	srv := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.Port)

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
