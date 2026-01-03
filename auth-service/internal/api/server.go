package api

import (
	"context"
	"fmt"
	"net/http"
	"github.com/portfolio/auth-service/internal/api/handlers"
	"github.com/portfolio/auth-service/internal/api/middleware"
	"github.com/portfolio/auth-service/internal/config"
	"github.com/portfolio/auth-service/internal/repository"
	"github.com/portfolio/auth-service/internal/service"
	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	config     *config.Config
	logger     *zap.Logger
	router     *gin.Engine
	db         *gorm.DB
	httpServer *http.Server
}

func NewServer(cfg *config.Config, zapLogger *zap.Logger) *Server {
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(
		userRepo,
		cfg.JWT.Secret,
		int(cfg.JWT.AccessExpiry.Minutes()),
		int(cfg.JWT.RefreshExpiry.Hours()),
	)

	authHandler := handlers.NewAuthHandler(authService)

	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.Logger(zapLogger))
	router.Use(middleware.Recovery(zapLogger))

	router.GET("/healthz", func(c *gin.Context) {
		c.String(200, "ok")
	})

	v1 := router.Group("/api/v1/auth")
	{
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)
		v1.POST("/refresh", authHandler.Refresh)
		v1.POST("/verify", authHandler.Verify)
		v1.POST("/logout", authHandler.Logout)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: router,
	}

	return &Server{
		config:     cfg,
		logger:     zapLogger,
		router:     router,
		db:         db,
		httpServer: httpServer,
	}
}

func (s *Server) Start() error {
	s.logger.Info("Starting auth service",
		zap.String("host", s.config.Server.Host),
		zap.String("port", s.config.Server.Port),
	)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

