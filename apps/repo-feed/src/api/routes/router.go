package routes

import (
	"repo_feed/src/api/controllers"
	"repo_feed/src/config"
	"repo_feed/src/domain/services"

	authMiddleware "github.com/TheRayquaza/newsbro/apps/libs/auth/middleware"
	middleware "github.com/TheRayquaza/newsbro/apps/libs/utils/middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "repo_feed/docs"

	"github.com/gin-gonic/gin"
	"log"
)

func SetupRouter(cfg *config.Config, feedService *services.FeedService, rescoringService *services.FeedRescoringService) *gin.Engine {
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
	feedController := controllers.NewFeedController(feedService, cfg.DefaultModel)
	//rescoringController := controllers.NewRescoringController(rescoringService)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		ctx := c.Request.Context()
		if err := feedService.RDB.Ping(ctx).Err(); err != nil {
			log.Println("Health check failed:", err)
			c.JSON(500, gin.H{"status": "unhealthy"})
			return
		}
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Readiness check
	router.GET("/ready", func(c *gin.Context) {
		ctx := c.Request.Context()
		if err := feedService.RDB.Ping(ctx).Err(); err != nil {
			log.Println("Readiness check failed:", err)
			c.JSON(500, gin.H{"status": "unready"})
			return
		}
		c.JSON(200, gin.H{"status": "ready"})
	})

	// Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		protected := v1.Group("/")
		protected.Use(authMiddleware.AuthMiddleware(cfg.JWTSecret, cfg.LoginRedirectURL))
		{
			// User routes
			feeds := protected.Group("/feed")
			{
				feeds.GET("", feedController.GetUserFeed)
				/*
					rescoring := feeds.Group("/rescoring")
					rescoring.Use(authMiddleware.AdminMiddleware())
					{
						rescoring.GET("/status", rescoringController.GetRescoringStatus)
						rescoring.POST("/reset", rescoringController.ResetRescoringCursor)
					}
				*/
			}
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
