package main

import (
	"context"
	"encoding/json"
	"log"
	"github.com/portfolio/backend/internal/config"
	"github.com/portfolio/backend/internal/model"
	"github.com/portfolio/backend/internal/repository"
	"time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/datatypes"
	"github.com/google/uuid"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	ctx := context.Background()

	// Initialize repositories
	articleRepo := repository.NewArticleRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	portfolioRepo := repository.NewPortfolioRepository(db)

	// Seed portfolio
	if err := seedPortfolio(ctx, portfolioRepo, cfg); err != nil {
		log.Printf("Error seeding portfolio: %v", err)
	}

	// Seed articles
	if err := seedArticles(ctx, articleRepo); err != nil {
		log.Printf("Error seeding articles: %v", err)
	}

	// Seed projects
	if err := seedProjects(ctx, projectRepo); err != nil {
		log.Printf("Error seeding projects: %v", err)
	}

	log.Println("Seeding completed successfully!")
}

func seedPortfolio(ctx context.Context, repo repository.PortfolioRepository, cfg *config.Config) error {
	socialLinksJSON, _ := json.Marshal(map[string]interface{}{
		"github":   cfg.Seeder.PortfolioGithub,
		"linkedin": cfg.Seeder.PortfolioLinkedin,
		"phone":    cfg.Seeder.PortfolioPhone,
	})
	
	settingsJSON, _ := json.Marshal(map[string]interface{}{
		"theme":    cfg.Seeder.PortfolioTheme,
		"language": cfg.Seeder.PortfolioLanguage,
	})

	portfolio := &model.Portfolio{
		Name:        cfg.Seeder.PortfolioName,
		Title:       cfg.Seeder.PortfolioTitle,
		Bio:         cfg.Seeder.PortfolioBio,
		Email:       cfg.Seeder.PortfolioEmail,
		SocialLinks: datatypes.JSON(socialLinksJSON),
		Settings:    datatypes.JSON(settingsJSON),
	}

	// Use CreateOrUpdate which handles both create and update
	return repo.CreateOrUpdate(ctx, portfolio)
}

func seedArticles(ctx context.Context, repo repository.ArticleRepository) error {
	authorID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	now := time.Now()

	articles := []*model.Article{
		{
			Title:      "Building Scalable Real-Time Bidding Systems with Go and Kafka",
			Slug:       "real-time-bidding-go-kafka",
			Excerpt:    "Learn how to architect high-performance auction-based advertising systems handling millions of bid requests with sub-100ms latency using Go and Kafka event streaming.",
			Content:    `# Building Scalable Real-Time Bidding Systems

Real-time bidding (RTB) systems require ultra-low latency and high throughput. In this article, I'll share my experience building such systems at Gowit Technology.

## Architecture Overview

The system processes real-time bid requests using:
- **Go** for high-performance request handling
- **Kafka** for event streaming and async processing
- **Redis** for caching and rate limiting
- **ClickHouse** for analytics and metrics

## Key Challenges

1. **Latency**: Sub-100ms response time requirement
2. **Throughput**: Handling millions of requests per second
3. **Reliability**: Zero-downtime deployments
4. **Scalability**: Horizontal scaling capabilities

## Solution

We implemented a microservices architecture with:
- Stateless API servers
- Event-driven processing
- Intelligent caching strategies
- Kubernetes orchestration

The result? A system that handles 10M+ requests per day with 99.9% uptime.`,
			AuthorID:   authorID,
			Published:  true,
			PublishedAt: &now,
		},
		{
			Title:      "Kubernetes Microservices: Best Practices from Production",
			Slug:       "kubernetes-microservices-best-practices",
			Excerpt:    "Production-tested strategies for deploying and managing microservices on Kubernetes, including service mesh, monitoring, and scaling patterns.",
			Content:    `# Kubernetes Microservices: Best Practices

After deploying multiple microservices to production, I've learned valuable lessons about Kubernetes orchestration.

## Service Design

- **Stateless services**: Enable easy scaling
- **Health checks**: Liveness and readiness probes
- **Resource limits**: Prevent resource starvation
- **ConfigMaps & Secrets**: Externalize configuration

## Monitoring & Observability

- Prometheus for metrics
- Grafana for visualization
- Structured logging with ELK stack
- Distributed tracing

## Scaling Strategies

- Horizontal Pod Autoscaling (HPA)
- Vertical Pod Autoscaling (VPA)
- Cluster Autoscaling

These practices ensure reliable, scalable microservices deployments.`,
			AuthorID:   authorID,
			Published:  true,
			PublishedAt: &now,
		},
		{
			Title:      "React Performance Optimization: Real-World Techniques",
			Slug:       "react-performance-optimization",
			Excerpt:    "Practical techniques for optimizing React applications, including code splitting, memoization, and virtual DOM optimization strategies.",
			Content:    `# React Performance Optimization

Building performant React applications requires understanding rendering behavior and optimization techniques.

## Code Splitting

Use React.lazy() and Suspense for route-based code splitting.

## Memoization

- React.memo() for component memoization
- useMemo() for expensive computations
- useCallback() for function references

## Virtual DOM Optimization

- Key props for list items
- Avoid inline object creation
- Minimize re-renders

These techniques can improve performance by 50%+.`,
			AuthorID:   authorID,
			Published:  true,
			PublishedAt: &now,
		},
		{
			Title:      "Event-Driven Architecture with Kafka: A Practical Guide",
			Slug:       "event-driven-architecture-kafka",
			Excerpt:    "Implementing event-driven microservices using Kafka, covering producer patterns, consumer groups, and event sourcing strategies.",
			Content:    `# Event-Driven Architecture with Kafka

Event-driven architecture enables loose coupling and scalability. Here's how we implemented it.

## Kafka Topics & Partitions

- Design topics for domain boundaries
- Use partitioning for parallel processing
- Consider replication factor for availability

## Producer Patterns

- Idempotent producers
- Batch processing
- Compression for efficiency

## Consumer Groups

- Parallel processing with consumer groups
- Offset management
- Error handling and retries

This architecture supports our real-time advertising platform processing millions of events daily.`,
			AuthorID:   authorID,
			Published:  true,
			PublishedAt: &now,
		},
		{
			Title:      "Redis Caching Strategies for High-Performance APIs",
			Slug:       "redis-caching-strategies",
			Excerpt:    "Effective caching patterns using Redis to reduce database load and improve API response times in high-traffic applications.",
			Content:    `# Redis Caching Strategies

Caching is crucial for high-performance APIs. Here are proven strategies.

## Cache Patterns

1. **Cache-Aside**: Application manages cache
2. **Write-Through**: Write to cache and DB
3. **Write-Back**: Write to cache, async to DB

## TTL Management

- Short TTL for frequently changing data
- Long TTL for stable data
- Cache invalidation on updates

## Key Design

- Consistent naming conventions
- Namespace prefixes
- Versioning for schema changes

These strategies reduced our API latency by 70%.`,
			AuthorID:   authorID,
			Published:  true,
			PublishedAt: &now,
		},
	}

	for _, article := range articles {
		existing, _ := repo.GetBySlug(ctx, article.Slug)
		if existing == nil {
			if err := repo.Create(ctx, article); err != nil {
				log.Printf("Error creating article %s: %v", article.Slug, err)
			}
		}
	}

	return nil
}

func seedProjects(ctx context.Context, repo repository.ProjectRepository) error {
	projects := []*model.Project{
		{
			Name:        "Real-Time Bidding Platform",
			Description: "Enterprise-scale digital advertising platform with real-time bidding and multi-channel campaign management. Handles high-volume bid requests with sub-100ms latency using Go and Kafka.",
			GithubURL:   "https://github.com/kmuhsinn/rtb-platform",
			LiveURL:     "",
			Technologies: model.StringArray{"Go", "Kubernetes", "React", "TypeScript", "Kafka", "Redis", "ClickHouse", "ScyllaDB", "PostgreSQL"},
			Featured:    true,
		},
		{
			Name:        "Campaign Management Dashboard",
			Description: "Comprehensive advertiser and retailer dashboard built with React and TypeScript. Enables self-service ad creation, placement management, and real-time performance analytics.",
			GithubURL:   "https://github.com/kmuhsinn/campaign-dashboard",
			LiveURL:     "",
			Technologies: model.StringArray{"React", "TypeScript", "Go", "PostgreSQL", "Redis"},
			Featured:    true,
		},
		{
			Name:        "Data Pipeline System",
			Description: "Scalable data pipeline processing real-time sales and impression data, computing key advertising metrics (CTR, ROAS, conversion rates) using ClickHouse and Redis.",
			GithubURL:   "https://github.com/kmuhsinn/data-pipeline",
			LiveURL:     "",
			Technologies: model.StringArray{"Go", "Python", "ClickHouse", "Redis", "Kafka", "Kubernetes"},
			Featured:    true,
		},
		{
			Name:        "SOAR Platform",
			Description: "Security Orchestration, Automation, and Response (SOAR) platform with advanced threat detection and automated incident response. Built with FastAPI, Elasticsearch, and microservices architecture.",
			GithubURL:   "https://github.com/kmuhsinn/soar-platform",
			LiveURL:     "",
			Technologies: model.StringArray{"Python", "FastAPI", "React", "GraphQL", "Elasticsearch", "MongoDB", "Docker", "Kubernetes"},
			Featured:    false,
		},
		{
			Name:        "UAV Control System",
			Description: "AI-powered UAV control systems and data processing platforms for autonomous flight operations. Integrates computer vision and deep learning models on edge devices.",
			GithubURL:   "https://github.com/kmuhsinn/uav-control",
			LiveURL:     "",
			Technologies: model.StringArray{"Python", "C++", "TypeScript", "gRPC", "MAVLink", "Computer Vision", "Deep Learning", "Edge Computing"},
			Featured:    false,
		},
		{
			Name:        "Energy Management Platform",
			Description: "Real-time energy management platform for solar power facilities. Features dynamic pricing engine and real-time data ingestion from multiple solar plants.",
			GithubURL:   "https://github.com/kmuhsinn/energy-platform",
			LiveURL:     "",
			Technologies: model.StringArray{"Go", "PHP", "React", "PostgreSQL", "MySQL", "Docker", "Kubernetes"},
			Featured:    false,
		},
		{
			Name:        "OpenAPI Application Generator",
			Description: "OpenAPI application generator that auto-generates fully functional APIs from Swagger documentation, reducing API development time by 80%.",
			GithubURL:   "https://github.com/kmuhsinn/openapi-generator",
			LiveURL:     "",
			Technologies: model.StringArray{"Python", "FastAPI", "OpenAPI", "Swagger"},
			Featured:    false,
		},
	}

	for _, project := range projects {
		// Try to create - if exists, will fail but that's ok
		if err := repo.Create(ctx, project); err != nil {
			// Project might already exist, log but continue
			log.Printf("Project %s might already exist: %v", project.Name, err)
		} else {
			log.Printf("Created project: %s", project.Name)
		}
	}

	return nil
}

