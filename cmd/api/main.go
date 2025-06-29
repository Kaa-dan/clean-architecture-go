package main

import (
	"github.com/kaa-dan/clean-architecture-go/internal/config"
	"github.com/kaa-dan/clean-architecture-go/pkg/logger"
)

func main() {
	// Load configuration

	cfg := config.Load()

	// Initialize logger
	logger.Init(cfg.LogLevel)
}
