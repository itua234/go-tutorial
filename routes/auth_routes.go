package routes

import (
	controllers "confam-api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/register", controllers.Register)
			auth.POST("/password/forgot", controllers.ForgotPassword)
			auth.POST("/password/reset", controllers.PasswordReset)
			auth.PUT("/password/change", controllers.ChangePassword)
			//auth.POST("/token/refresh", controllers.RefreshToken)
		}
	}
}
