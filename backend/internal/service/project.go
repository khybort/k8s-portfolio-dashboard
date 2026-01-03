package service

import (
	"context"
	"fmt"
	"github.com/portfolio/backend/internal/cache"
	"github.com/portfolio/backend/internal/kafka"
	"github.com/portfolio/backend/internal/model"
	"github.com/portfolio/backend/internal/repository"
)

type ProjectService interface {
	GetProjects(ctx context.Context, page, limit int, featured *bool) ([]model.Project, int64, error)
	GetProjectByID(ctx context.Context, id string) (*model.Project, error)
	CreateProject(ctx context.Context, project *model.Project) error
	UpdateProject(ctx context.Context, id string, project *model.Project) error
	DeleteProject(ctx context.Context, id string) error
}

type projectService struct {
	repo  repository.ProjectRepository
	kafka *kafka.Producer
	cache *cache.RedisCache
}

func NewProjectService(repo repository.ProjectRepository, kafka *kafka.Producer, cache *cache.RedisCache) ProjectService {
	return &projectService{
		repo:  repo,
		kafka: kafka,
		cache: cache,
	}
}

func (s *projectService) GetProjects(ctx context.Context, page, limit int, featured *bool) ([]model.Project, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.List(ctx, page, limit, featured)
}

func (s *projectService) GetProjectByID(ctx context.Context, id string) (*model.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *projectService) CreateProject(ctx context.Context, project *model.Project) error {
	if err := s.repo.Create(ctx, project); err != nil {
		return err
	}

	// Publish Kafka event
	s.kafka.PublishProjectCreated(ctx, project)

	// Invalidate cache
	s.cache.DeletePattern(ctx, "projects:*")

	return nil
}

func (s *projectService) UpdateProject(ctx context.Context, id string, project *model.Project) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err // Repository already returns ErrProjectNotFound
	}

	project.ID = existing.ID
	if err := s.repo.Update(ctx, project); err != nil {
		return err
	}

	// Publish Kafka event
	s.kafka.PublishProjectUpdated(ctx, project)

	// Invalidate cache
	s.cache.DeletePattern(ctx, "projects:*")
	s.cache.Delete(ctx, fmt.Sprintf("project:detail:%s", id))

	return nil
}

func (s *projectService) DeleteProject(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Publish Kafka event
	s.kafka.PublishProjectDeleted(ctx, id)

	// Invalidate cache
	s.cache.DeletePattern(ctx, "projects:*")
	s.cache.Delete(ctx, fmt.Sprintf("project:detail:%s", id))

	return nil
}

