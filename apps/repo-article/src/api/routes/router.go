package routes

import (
	"repo_article/src/api/controllers"
	"repo_article/src/config"
	"repo_article/src/domain/services"

	authMiddleware "github.com/TheRayquaza/newsbro/apps/libs/auth/middleware"
	middleware "github.com/TheRayquaza/newsbro/apps/libs/utils/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "repo_article/docs"
)

func SetupRouter(cfg *config.Config, articleService *services.ArticleService, feedbackService *services.FeedbackService) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	if cfg.Environment != "dev" {
		router.Use(middleware.CORSMiddleware(cfg.FrontendOrigin))
	}
	router.Use(middleware.Secure())

	// Controllers
	articleController := controllers.NewArticleController(articleService)
	feedbackController := controllers.NewFeedbackController(feedbackService)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		if err := articleService.Db.Raw("SELECT 1").Error; err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "reason": err})
			return
		}
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Readiness check
	router.GET("/ready", func(c *gin.Context) {
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
		protected.Use(authMiddleware.AuthMiddleware(cfg.JWTSecret, cfg.LoginRedirectURL))
		{
			// Article routes
			articles := protected.Group("/articles")
			{
				articles.GET("", articleController.GetArticles)
				articles.GET("/categories", articleController.GetCategories)
				articles.GET("/subcategories", articleController.GetSubcategories)
				articles.GET("/history", articleController.GetArticleHistory)
				articles.POST("/ingestion", authMiddleware.AdminMiddleware(), articleController.TriggerArticleIngestion)
				articles.GET("/:id", articleController.GetArticle)
				articles.POST("", authMiddleware.AdminMiddleware(), articleController.CreateArticle)
				articles.PUT("/:id", authMiddleware.AdminMiddleware(), articleController.UpdateArticle)
				articles.DELETE("/:id", authMiddleware.AdminMiddleware(), articleController.DeleteArticle)

				// Article feedback routes
				articles.GET("/:id/feedback", authMiddleware.AdminMiddleware(), feedbackController.GetArticleFeedback)
				articles.POST("/:id/feedback", feedbackController.CreateFeedback)
				articles.DELETE("/:id/feedback", feedbackController.DeleteFeedback)
			}

			// Feedback routes
			feedback := protected.Group("/feedback")
			{
				// Get user's feedback history
				feedback.GET("/my", feedbackController.GetUserFeedback)

				// Admin routes for feedback management
				feedback.GET("/csv", authMiddleware.AdminMiddleware(), feedbackController.ExportFeedbackCSV)
				feedback.GET("/stats", authMiddleware.AdminMiddleware(), feedbackController.GetFeedbackStats)
				feedback.GET("/all", authMiddleware.AdminMiddleware(), feedbackController.GetAllFeedback)
				feedback.POST("/ingest", authMiddleware.AdminMiddleware(), feedbackController.TriggerIngestFeedback)
			}
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return router
}
