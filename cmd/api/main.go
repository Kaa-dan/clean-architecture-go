package main

import (
	"context"
	"log"

	"github.com/kaa-dan/clean-architecture-go/internal/config"
	"github.com/kaa-dan/clean-architecture-go/internal/infrastructure/database"
	"github.com/kaa-dan/clean-architecture-go/internal/infrastructure/repositories"
	"github.com/kaa-dan/clean-architecture-go/internal/infrastructure/security"
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
}
