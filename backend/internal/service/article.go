package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/portfolio/backend/internal/cache"
	"github.com/portfolio/backend/internal/kafka"
	"github.com/portfolio/backend/internal/model"
	"github.com/portfolio/backend/internal/repository"
	"time"
)

type ArticleService interface {
	GetArticles(ctx context.Context, page, limit int) ([]model.Article, int64, error)
	GetArticleByID(ctx context.Context, id string) (*model.Article, error)
	GetArticleBySlug(ctx context.Context, slug string) (*model.Article, error)
	CreateArticle(ctx context.Context, article *model.Article) error
	UpdateArticle(ctx context.Context, id string, article *model.Article) error
	DeleteArticle(ctx context.Context, id string) error
}

type articleService struct {
	repo   repository.ArticleRepository
	kafka  *kafka.Producer
	cache  *cache.RedisCache
}

func NewArticleService(repo repository.ArticleRepository, kafka *kafka.Producer, cache *cache.RedisCache) ArticleService {
	return &articleService{
		repo:  repo,
		kafka: kafka,
		cache: cache,
	}
}

func (s *articleService) GetArticles(ctx context.Context, page, limit int) ([]model.Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.List(ctx, page, limit, true)
}

func (s *articleService) GetArticleByID(ctx context.Context, id string) (*model.Article, error) {
	// Try cache first
	cached, err := s.cache.GetArticle(ctx, id)
	if err == nil {
		var article model.Article
		json.Unmarshal(cached, &article)
		return &article, nil
	}

	// Cache miss - get from database
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store in cache
	s.cache.SetArticle(ctx, id, article, 10*time.Minute)

	return article, nil
}

func (s *articleService) GetArticleBySlug(ctx context.Context, slug string) (*model.Article, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *articleService) CreateArticle(ctx context.Context, article *model.Article) error {
	if article.Published {
		now := time.Now()
		article.PublishedAt = &now
	}

	if err := s.repo.Create(ctx, article); err != nil {
		return err
	}

	// Publish Kafka event
	s.kafka.PublishArticleCreated(ctx, article)

	// Invalidate cache
	s.cache.InvalidateArticles(ctx)

	return nil
}

func (s *articleService) UpdateArticle(ctx context.Context, id string, article *model.Article) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err // Repository already returns ErrArticleNotFound
	}

	article.ID = existing.ID
	if article.Published && existing.PublishedAt == nil {
		now := time.Now()
		article.PublishedAt = &now
	}

	if err := s.repo.Update(ctx, article); err != nil {
		return err
	}

	// Publish Kafka event
	s.kafka.PublishArticleUpdated(ctx, article)

	// Invalidate cache
	s.cache.InvalidateArticles(ctx)
	s.cache.Delete(ctx, fmt.Sprintf("article:detail:%s", id))

	return nil
}

func (s *articleService) DeleteArticle(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Publish Kafka event
	s.kafka.PublishArticleDeleted(ctx, id)

	// Invalidate cache
	s.cache.InvalidateArticles(ctx)
	s.cache.Delete(ctx, fmt.Sprintf("article:detail:%s", id))

	return nil
}

