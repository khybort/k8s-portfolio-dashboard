package handlers

import (
	"errors"
	"net/http"
	"github.com/portfolio/backend/internal/model"
	"github.com/portfolio/backend/internal/repository"
	"github.com/portfolio/backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	service service.ProjectService
}

func NewProjectHandler(service service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) GetProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	var featured *bool
	if featuredStr := c.Query("featured"); featuredStr != "" {
		feat := featuredStr == "true"
		featured = &feat
	}

	projects, total, err := h.service.GetProjects(c.Request.Context(), page, limit, featured)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (int(total) + limit - 1) / limit
	c.JSON(http.StatusOK, gin.H{
		"data": projects,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": totalPages,
		},
	})
}

func (h *ProjectHandler) GetProjectByID(c *gin.Context) {
	id := c.Param("id")
	project, err := h.service.GetProjectByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var project struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		GithubURL   string   `json:"github_url"`
		LiveURL     string   `json:"live_url"`
		Technologies []string `json:"technologies"`
		Featured    bool     `json:"featured"`
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectModel := &model.Project{
		Name:        project.Name,
		Description: project.Description,
		GithubURL:   project.GithubURL,
		LiveURL:     project.LiveURL,
		Technologies: model.StringArray(project.Technologies),
		Featured:    project.Featured,
	}

	if err := h.service.CreateProject(c.Request.Context(), projectModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, projectModel)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	
	var project struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		GithubURL   string   `json:"github_url"`
		LiveURL     string   `json:"live_url"`
		Technologies []string `json:"technologies"`
		Featured    bool     `json:"featured"`
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectModel := &model.Project{
		Name:        project.Name,
		Description: project.Description,
		GithubURL:   project.GithubURL,
		LiveURL:     project.LiveURL,
		Technologies: model.StringArray(project.Technologies),
		Featured:    project.Featured,
	}

	if err := h.service.UpdateProject(c.Request.Context(), id, projectModel); err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projectModel)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteProject(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

