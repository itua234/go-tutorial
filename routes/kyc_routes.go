package routes

import (
	controllers "confam-api/controllers"
	database "confam-api/database"
	"confam-api/middlewares"
	repositories "confam-api/repositories"
	services "confam-api/services"

	"github.com/gin-gonic/gin"
)

func RegisterKycRoutes(router *gin.Engine) {
	db := database.DB

	//companyRepo := repositories.NewCompanyRepository(db)
	customerRepo := repositories.NewCustomerRepository(db)
	requestRepo := repositories.NewRequestRepository(db)

	// 3. Create Service instances, injecting repositories
	kycService := services.NewKYCService(customerRepo, requestRepo)

	// 4. Create Controller instances, injecting services
	kycController := controllers.NewKycController(kycService)

	api := router.Group("/api/v1")
	{
		kyc := api.Group("/allow")
		{
			kyc.POST(
				"/",
				middlewares.AuthenticateAppBySecretKey(database.DB),
				kycController.InitiateKyc,
			)
			//kyc.GET("/:kyc_token", kycController.FetchKycRequest)
		}
	}
}
