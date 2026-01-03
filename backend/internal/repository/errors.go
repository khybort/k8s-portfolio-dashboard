package repository

import "errors"

// Common repository errors
var (
	ErrArticleNotFound   = errors.New("article not found")
	ErrProjectNotFound   = errors.New("project not found")
	ErrPortfolioNotFound = errors.New("portfolio not found")
)

