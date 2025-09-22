package routes

import (
	"repo_account/src/api/controllers"
	"repo_account/src/api/middleware"
	"repo_account/src/config"
	"repo_account/src/domain/services"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "repo_account/docs"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, authService *services.AuthService, userService *services.UserService) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.Secure())
	router.Use(middleware.RemoveTrailingSlash())

	// Controllers
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		if err := userService.Db.Raw("SELECT 1").Error; err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "reason": err})
			return
		}
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Readiness check
	router.GET("/readiness", func(c *gin.Context) {
		if err := userService.Db.Raw("SELECT 1").Error; err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "reason": err})
			return
		}
		c.JSON(200, gin.H{"status": "ready"})
	})

	// OAuth routes for EPITA
	router.GET("/complete/epita", authController.OAuthCallback)

	// Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", middleware.AuthMiddleware(authService), authController.RefreshToken)
			auth.GET("/oauth/login", authController.OAuthLogin)
			auth.GET("/callback", authController.OAuthCallback)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", userController.GetProfile)
				users.PUT("/profile", userController.UpdateProfile)
				//users.DELETE("/profile", userController.DeleteProfile)
				users.GET("/", userController.GetUsers) // Admin only in production
			}
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
