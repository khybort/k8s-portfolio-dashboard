package repository

import (
	"context"
	"errors"
	"github.com/portfolio/backend/internal/model"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *model.Project) error
	GetByID(ctx context.Context, id string) (*model.Project, error)
	List(ctx context.Context, page, limit int, featured *bool) ([]model.Project, int64, error)
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, id string) error
	WithTransaction(ctx context.Context, fn func(*gorm.DB) error) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) WithTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

func (r *projectRepository) Create(ctx context.Context, project *model.Project) error {
	if err := r.db.WithContext(ctx).Create(project).Error; err != nil {
		return err
	}
	return nil
}

func (r *projectRepository) GetByID(ctx context.Context, id string) (*model.Project, error) {
	var project model.Project
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&project).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) List(ctx context.Context, page, limit int, featured *bool) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Project{})
	
	if featured != nil {
		query = query.Where("featured = ?", *featured)
	}

	// Count total (before pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Get paginated results with optimized query
	offset := (page - 1) * limit
	err := query.
		Select("id", "name", "description", "github_url", "live_url", "technologies", "featured", "created_at", "updated_at").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&projects).Error
	
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *projectRepository) Update(ctx context.Context, project *model.Project) error {
	// Use Updates to only update non-zero fields
	result := r.db.WithContext(ctx).
		Model(project).
		Where("id = ?", project.ID).
		Updates(map[string]interface{}{
			"name":        project.Name,
			"description": project.Description,
			"github_url":  project.GithubURL,
			"live_url":    project.LiveURL,
			"technologies": project.Technologies,
			"featured":    project.Featured,
		})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return ErrProjectNotFound
	}
	
	return nil
}

func (r *projectRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Project{})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return ErrProjectNotFound
	}
	
	return nil
}

