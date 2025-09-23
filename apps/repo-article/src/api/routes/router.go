package routes

import (
	"repo_article/src/api/controllers"
	"repo_article/src/api/middleware"
	"repo_article/src/config"
	"repo_article/src/domain/services"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "repo_article/docs"
)

func SetupRouter(cfg *config.Config, articleService *services.ArticleService, authService *services.AuthService) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.Secure())
	router.Use(middleware.RemoveTrailingSlash())

	// Controllers
	articleController := controllers.NewArticleController(articleService)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		if err := articleService.Db.Raw("SELECT 1").Error; err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "reason": err})
			return
		}
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Readiness check
	router.GET("/readiness", func(c *gin.Context) {
		if err := articleService.Db.Raw("SELECT 1").Error; err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "reason": err})
			return
		}
		c.JSON(200, gin.H{"status": "ready"})
	})

	// Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// Article routes
			articles := protected.Group("/articles")
			{
				articles.POST("", articleController.CreateArticle)
				articles.GET("", articleController.GetArticles)
				articles.GET("/categories", articleController.GetCategories)
				articles.GET("/subcategories", articleController.GetSubcategories)
				articles.GET("/:id", articleController.GetArticle)
				articles.PUT("/:id", articleController.UpdateArticle)
				articles.DELETE("/:id", articleController.DeleteArticle)
			}
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return router
}
