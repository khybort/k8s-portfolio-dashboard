package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"github.com/portfolio/backend/internal/api/handlers"
	"github.com/portfolio/backend/internal/api/middleware"
	"github.com/portfolio/backend/internal/config"
	"github.com/portfolio/backend/internal/model"
	"github.com/portfolio/backend/internal/repository"
	"github.com/portfolio/backend/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/portfolio/backend/internal/cache"
	"github.com/portfolio/backend/internal/kafka"
)

type Server struct {
	config   *config.Config
	logger   *zap.Logger
	router   *gin.Engine
	db       *gorm.DB
	httpServer *http.Server
}

func NewServer(cfg *config.Config, zapLogger *zap.Logger) *Server {
	// Initialize database
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		PrepareStmt: true, // Use prepared statements for better performance
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto migrate models
	if err := autoMigrate(db, zapLogger); err != nil {
		zapLogger.Fatal("Failed to auto migrate database", zap.Error(err))
	}

	// Initialize Redis cache
	redisCache := cache.NewRedisCache(
		fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		cfg.Redis.Password,
		cfg.Redis.DB,
	)

	// Initialize Kafka producer
	kafkaProducer := kafka.NewProducer(cfg.Kafka.Brokers)
	defer kafkaProducer.Close()

	// Initialize repositories
	articleRepo := repository.NewArticleRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	portfolioRepo := repository.NewPortfolioRepository(db)

	// Initialize services (with Kafka and Redis)
	articleService := service.NewArticleService(articleRepo, kafkaProducer, redisCache)
	projectService := service.NewProjectService(projectRepo, kafkaProducer, redisCache)
	portfolioService := service.NewPortfolioService(portfolioRepo)

	// Initialize handlers
	articleHandler := handlers.NewArticleHandler(articleService)
	projectHandler := handlers.NewProjectHandler(projectService)
	portfolioHandler := handlers.NewPortfolioHandler(portfolioService)

	// Setup router
	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.Logger(zapLogger))
	router.Use(middleware.Recovery(zapLogger))

	// Health check
	router.GET("/healthz", func(c *gin.Context) {
		c.String(200, "ok")
	})

	// Public API routes
	v1 := router.Group("/api/v1")
	{
		// Articles
		v1.GET("/articles", articleHandler.GetArticles)
		v1.GET("/articles/:id", articleHandler.GetArticleByID)
		v1.GET("/articles/slug/:slug", articleHandler.GetArticleBySlug)

		// Projects
		v1.GET("/projects", projectHandler.GetProjects)
		v1.GET("/projects/:id", projectHandler.GetProjectByID)

		// Portfolio
		v1.GET("/portfolio", portfolioHandler.GetPortfolio)
	}

	// Admin API routes (require authentication)
	admin := v1.Group("/admin")
	admin.Use(middleware.Auth(cfg.Auth.ServiceURL))
	{
		// Articles
		admin.POST("/articles", articleHandler.CreateArticle)
		admin.PUT("/articles/:id", articleHandler.UpdateArticle)
		admin.DELETE("/articles/:id", articleHandler.DeleteArticle)

		// Projects
		admin.POST("/projects", projectHandler.CreateProject)
		admin.PUT("/projects/:id", projectHandler.UpdateProject)
		admin.DELETE("/projects/:id", projectHandler.DeleteProject)

		// Portfolio
		admin.PUT("/portfolio", portfolioHandler.UpdatePortfolio)
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
	s.logger.Info("Starting server", 
		zap.String("host", s.config.Server.Host),
		zap.String("port", s.config.Server.Port),
	)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	// Close database connection
	sqlDB, err := s.db.DB()
	if err == nil {
		sqlDB.Close()
	}
	return s.httpServer.Shutdown(ctx)
}

// autoMigrate runs GORM AutoMigrate for all models
func autoMigrate(db *gorm.DB, logger *zap.Logger) error {
	models := []interface{}{
		&model.Article{},
		&model.Project{},
		&model.Portfolio{},
	}

	for _, m := range models {
		if err := db.AutoMigrate(m); err != nil {
			logger.Error("Failed to auto migrate model", zap.Error(err), zap.String("model", fmt.Sprintf("%T", m)))
			return err
		}
	}

	logger.Info("Database auto migration completed successfully")
	return nil
}

