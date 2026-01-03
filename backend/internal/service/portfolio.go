package service

import (
	"context"
	"errors"
	"github.com/portfolio/backend/internal/model"
	"github.com/portfolio/backend/internal/repository"
)

type PortfolioService interface {
	GetPortfolio(ctx context.Context) (*model.Portfolio, error)
	CreateOrUpdatePortfolio(ctx context.Context, portfolio *model.Portfolio) error
	UpdatePortfolio(ctx context.Context, portfolio *model.Portfolio) error
}

type portfolioService struct {
	repo repository.PortfolioRepository
}

func NewPortfolioService(repo repository.PortfolioRepository) PortfolioService {
	return &portfolioService{repo: repo}
}

func (s *portfolioService) GetPortfolio(ctx context.Context) (*model.Portfolio, error) {
	return s.repo.Get(ctx)
}

func (s *portfolioService) CreateOrUpdatePortfolio(ctx context.Context, portfolio *model.Portfolio) error {
	return s.repo.CreateOrUpdate(ctx, portfolio)
}

func (s *portfolioService) UpdatePortfolio(ctx context.Context, portfolio *model.Portfolio) error {
	if err := s.repo.Update(ctx, portfolio); err != nil {
		if errors.Is(err, repository.ErrPortfolioNotFound) {
			// If not found, create it
			return s.repo.CreateOrUpdate(ctx, portfolio)
		}
		return err
	}
	return nil
}

