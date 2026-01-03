package repository

import (
	"context"
	"errors"
	"github.com/portfolio/backend/internal/model"
	"gorm.io/gorm"
)

type PortfolioRepository interface {
	Get(ctx context.Context) (*model.Portfolio, error)
	CreateOrUpdate(ctx context.Context, portfolio *model.Portfolio) error
	Update(ctx context.Context, portfolio *model.Portfolio) error
	WithTransaction(ctx context.Context, fn func(*gorm.DB) error) error
}

type portfolioRepository struct {
	db *gorm.DB
}

func NewPortfolioRepository(db *gorm.DB) PortfolioRepository {
	return &portfolioRepository{db: db}
}

func (r *portfolioRepository) WithTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

func (r *portfolioRepository) Get(ctx context.Context) (*model.Portfolio, error) {
	var portfolio model.Portfolio
	err := r.db.WithContext(ctx).First(&portfolio).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPortfolioNotFound
		}
		return nil, err
	}
	return &portfolio, nil
}

func (r *portfolioRepository) CreateOrUpdate(ctx context.Context, portfolio *model.Portfolio) error {
	// Try to get existing portfolio
	existing, err := r.Get(ctx)
	if err != nil && !errors.Is(err, ErrPortfolioNotFound) {
		return err
	}

	if existing != nil {
		// Update existing
		portfolio.ID = existing.ID
		return r.Update(ctx, portfolio)
	}

	// Create new
	return r.db.WithContext(ctx).Create(portfolio).Error
}

func (r *portfolioRepository) Update(ctx context.Context, portfolio *model.Portfolio) error {
	if portfolio.ID == (model.Portfolio{}).ID {
		return errors.New("portfolio ID is required for update")
	}

	result := r.db.WithContext(ctx).
		Model(portfolio).
		Where("id = ?", portfolio.ID).
		Updates(map[string]interface{}{
			"name":        portfolio.Name,
			"title":       portfolio.Title,
			"bio":         portfolio.Bio,
			"email":       portfolio.Email,
			"social_links": portfolio.SocialLinks,
			"settings":    portfolio.Settings,
		})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return ErrPortfolioNotFound
	}
	
	return nil
}

