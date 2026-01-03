package handlers

import (
	"errors"
	"net/http"
	"github.com/portfolio/backend/internal/model"
	"github.com/portfolio/backend/internal/repository"
	"github.com/portfolio/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type PortfolioHandler struct {
	service service.PortfolioService
}

func NewPortfolioHandler(service service.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{service: service}
}

func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
	portfolio, err := h.service.GetPortfolio(c.Request.Context())
	if err != nil {
		if errors.Is(err, repository.ErrPortfolioNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, portfolio)
}

func (h *PortfolioHandler) UpdatePortfolio(c *gin.Context) {
	var portfolio model.Portfolio

	if err := c.ShouldBindJSON(&portfolio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdatePortfolio(c.Request.Context(), &portfolio); err != nil {
		if errors.Is(err, repository.ErrPortfolioNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

