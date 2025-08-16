package routes

import (
	controllers "confam-api/controllers"
	database "confam-api/database"
	repositories "confam-api/repositories"
	services "confam-api/services"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine) {
	// This should be your GORM database instance
	db := database.DB

	// 2. Create Repository instances
	companyRepo := repositories.NewCompanyRepository(db)

	// 3. Create Service instances, injecting repositories
	authService := services.NewAuthService(companyRepo)

	// 4. Create Controller instances, injecting services
	authController := controllers.NewAuthController(authService)

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/register", controllers.Register)
			auth.POST("/password/forgot", controllers.ForgotPassword)
			auth.POST("/password/reset", controllers.PasswordReset)
			auth.PUT("/password/change", controllers.ChangePassword)
			//auth.POST("/token/refresh", controllers.RefreshToken)
		}
	}
}
