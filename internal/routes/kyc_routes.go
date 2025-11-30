package routes

import (
	controllers "confam-api/internal/controllers"
	database "confam-api/internal/database"
	"confam-api/internal/middlewares"
	client "confam-api/internal/redis"
	repositories "confam-api/internal/repositories"
	services "confam-api/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterKycRoutes(router *gin.Engine, rdb *client.Client) {
	db := database.DB

	//companyRepo := repositories.NewCompanyRepository(db)
	customerRepo := repositories.NewCustomerRepository(db)
	requestRepo := repositories.NewRequestRepository(db)
	//appRepo := repositories.NewAppRepository(db, rdb)

	// 3. Create Service instances, injecting repositories
	kycService := services.NewKYCService(customerRepo, requestRepo)
	webhookService := services.NewWebhookService()

	// 4. Create Controller instances, injecting services
	kycController := controllers.NewKycController(kycService, webhookService)

	api := router.Group("/api/v1")
	{
		kyc := api.Group("/allow")
		{
			kyc.POST(
				"/",
				middlewares.AuthenticateAppBySecretKey(database.DB, rdb),
				kycController.InitiateKyc,
			)
			kyc.GET("/:kyc_token", kycController.FetchKycRequest)
		}
	}
}
