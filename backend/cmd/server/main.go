package main

import (
	"log"
	"github.com/portfolio/backend/internal/api"
	"github.com/portfolio/backend/internal/config"
	"github.com/portfolio/backend/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	zapLogger, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	// Initialize and start server
	server := api.NewServer(cfg, zapLogger)
	if err := server.Start(); err != nil {
		zapLogger.Fatal("Failed to start server", zap.Error(err))
	}
}

