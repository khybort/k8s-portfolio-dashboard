package repository

import (
	"context"
	"errors"
	"github.com/portfolio/backend/internal/model"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *model.Article) error
	GetByID(ctx context.Context, id string) (*model.Article, error)
	GetBySlug(ctx context.Context, slug string) (*model.Article, error)
	List(ctx context.Context, page, limit int, published bool) ([]model.Article, int64, error)
	Update(ctx context.Context, article *model.Article) error
	Delete(ctx context.Context, id string) error
	WithTransaction(ctx context.Context, fn func(*gorm.DB) error) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) WithTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

func (r *articleRepository) Create(ctx context.Context, article *model.Article) error {
	if err := r.db.WithContext(ctx).Create(article).Error; err != nil {
		return err
	}
	return nil
}

func (r *articleRepository) GetByID(ctx context.Context, id string) (*model.Article, error) {
	var article model.Article
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&article).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) GetBySlug(ctx context.Context, slug string) (*model.Article, error) {
	var article model.Article
	err := r.db.WithContext(ctx).
		Where("slug = ? AND published = ?", slug, true).
		First(&article).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) List(ctx context.Context, page, limit int, published bool) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Article{})
	
	if published {
		query = query.Where("published = ?", true)
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
		Select("id", "title", "slug", "excerpt", "author_id", "published", "published_at", "created_at", "updated_at").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&articles).Error
	
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (r *articleRepository) Update(ctx context.Context, article *model.Article) error {
	// Use Updates to only update non-zero fields
	result := r.db.WithContext(ctx).
		Model(article).
		Where("id = ?", article.ID).
		Updates(map[string]interface{}{
			"title":       article.Title,
			"slug":        article.Slug,
			"excerpt":     article.Excerpt,
			"content":     article.Content,
			"published":   article.Published,
			"published_at": article.PublishedAt,
		})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return ErrArticleNotFound
	}
	
	return nil
}

func (r *articleRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Article{})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return ErrArticleNotFound
	}
	
	return nil
}

