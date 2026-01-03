package handlers

import (
	"errors"
	"net/http"
	"github.com/portfolio/backend/internal/model"
	"github.com/portfolio/backend/internal/repository"
	"github.com/portfolio/backend/internal/service"
	"strconv"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	service service.ArticleService
}

func NewArticleHandler(service service.ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

func (h *ArticleHandler) GetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	articles, total, err := h.service.GetArticles(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (int(total) + limit - 1) / limit
	c.JSON(http.StatusOK, gin.H{
		"data": articles,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": totalPages,
		},
	})
}

func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
	id := c.Param("id")
	article, err := h.service.GetArticleByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrArticleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) GetArticleBySlug(c *gin.Context) {
	slug := c.Param("slug")
	article, err := h.service.GetArticleBySlug(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, repository.ErrArticleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var article struct {
		Title     string `json:"title" binding:"required"`
		Slug      string `json:"slug" binding:"required"`
		Excerpt   string `json:"excerpt"`
		Content   string `json:"content" binding:"required"`
		Published bool   `json:"published"`
	}

	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get author_id from JWT token
	userID, _ := c.Get("user_id")
	userIDStr, _ := userID.(string)
	
	authorID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	articleModel := &model.Article{
		Title:     article.Title,
		Slug:      article.Slug,
		Excerpt:   article.Excerpt,
		Content:   article.Content,
		Published: article.Published,
		AuthorID:  authorID,
	}

	if err := h.service.CreateArticle(c.Request.Context(), articleModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, articleModel)
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	
	var article struct {
		Title     string `json:"title"`
		Slug      string `json:"slug"`
		Excerpt   string `json:"excerpt"`
		Content   string `json:"content"`
		Published bool   `json:"published"`
	}

	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	articleModel := &model.Article{
		Title:     article.Title,
		Slug:      article.Slug,
		Excerpt:   article.Excerpt,
		Content:   article.Content,
		Published: article.Published,
	}

	if err := h.service.UpdateArticle(c.Request.Context(), id, articleModel); err != nil {
		if errors.Is(err, repository.ErrArticleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articleModel)
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteArticle(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrArticleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

