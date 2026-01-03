package main

import (
	"log"
	"github.com/portfolio/auth-service/internal/api"
	"github.com/portfolio/auth-service/internal/config"
	"github.com/portfolio/auth-service/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	zapLogger, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	server := api.NewServer(cfg, zapLogger)
	if err := server.Start(); err != nil {
		zapLogger.Fatal("Failed to start server", zap.Error(err))
	}
}

