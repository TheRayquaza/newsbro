package controllers

import (
	"net/http"
	"strconv"

	"repo_article/src/api/dto"
	"repo_article/src/domain/entities"
	"repo_article/src/domain/services"
	//"repo_article/src/data/models"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	articleService *services.ArticleService
}

func NewArticleController(articleService *services.ArticleService) *ArticleController {
	return &ArticleController{
		articleService: articleService,
	}
}

// @Summary Create a new article
// @Description Create a new article with the provided data
// @Tags Articles
// @Accept json
// @Produce json
// @Param article body dto.ArticleCreateRequest true "Article data"
// @Success 201 {object} models.Article
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /articles [post]
func (ac *ArticleController) CreateArticle(c *gin.Context) {
	var req dto.ArticleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userID := user.(*entities.User).ID

	article, err := ac.articleService.CreateArticle(&req, userID)
	if err != nil {
		switch err.(type) {
		case *dto.ErrConflict:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case *dto.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, article)
}

// @Summary Get article by ID
// @Description Get a specific article by its ID
// @Tags Articles
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {object} models.Article
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /articles/{id} [get]
func (ac *ArticleController) GetArticle(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	article, err := ac.articleService.GetArticleByID(uint(id))
	if err != nil {
		if err.Error() == "article with id "+idParam+" not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}

// @Summary Get all articles
// @Description Get a list of articles with optional filtering and pagination
// @Tags Articles
// @Produce json
// @Param category query string false "Filter by category"
// @Param subcategory query string false "Filter by subcategory"
// @Param search query string false "Search in title and abstract"
// @Param limit query int false "Limit number of results" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /articles [get]
func (ac *ArticleController) GetArticles(c *gin.Context) {
	var filters dto.ArticleFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if filters.Limit <= 0 {
		filters.Limit = 10
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}

	articles, total, err := ac.articleService.GetArticles(&filters)
	if err != nil {
		switch err.(type) {
		case *dto.ErrBadRequest:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
		"total":    total,
		"limit":    filters.Limit,
		"offset":   filters.Offset,
	})
}

// @Summary Update article
// @Description Update an existing article by ID
// @Tags Articles
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Param article body dto.ArticleUpdateRequest true "Article update data"
// @Success 200 {object} models.Article
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /articles/{id} [put]
func (ac *ArticleController) UpdateArticle(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var req dto.ArticleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userID := user.(*entities.User).ID

	article, err := ac.articleService.UpdateArticle(uint(id), &req, userID)
	if err != nil {
		switch err.(type) {
		case *dto.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case *dto.ErrConflict:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, article)
}

// @Summary Delete article
// @Description Delete an article by ID
// @Tags Articles
// @Param id path int true "Article ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /articles/{id} [delete]
func (ac *ArticleController) DeleteArticle(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userID := user.(*entities.User).ID

	err = ac.articleService.DeleteArticle(uint(id), userID)
	if err != nil {
		switch err.(type) {
		case *dto.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get categories
// @Description Get all available article categories
// @Tags Articles
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /articles/categories [get]
func (ac *ArticleController) GetCategories(c *gin.Context) {
	categories, err := ac.articleService.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// @Summary Get subcategories
// @Description Get subcategories for a specific category or all subcategories
// @Tags Articles
// @Produce json
// @Param category query string false "Filter subcategories by category"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /articles/subcategories [get]
func (ac *ArticleController) GetSubcategories(c *gin.Context) {
	category := c.Query("category")

	subcategories, err := ac.articleService.GetSubcategories(category)
	if err != nil {
		switch err.(type) {
		case *dto.ErrBadRequest:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"subcategories": subcategories})
}
