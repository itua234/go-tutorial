package main

import (
	"fmt"
	"net/http"

	"confam-api/database"
	"confam-api/routes"

	"confam-api/middlewares"
	client "confam-api/utils"

	"confam-api/controllers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

// User structs for validation
type RegisterRequest struct {
	FirstName       string `json:"first_name" binding:"required,min=2,max=50" validate:"alpha"`
	LastName        string `json:"last_name" binding:"required,min=2,max=50" validate:"alpha"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8,max=100"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Phone           string `json:"phone" binding:"omitempty" validate:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Custom validator instance
var validate *validator.Validate

func main() {
	godotenv.Load()
	client.InitRedisClient() // Initialize Redis
	router := gin.Default()
	database.ConnectDatabase()

	router.MaxMultipartMemory = 8 << 20 //8Mb
	router.POST("/api/v1/upload", func(c *gin.Context) {
		// form, _ := c.MultipartForm
		// files := form.File["files"]
		file, _ := c.FormFile("file")
		log.Println(file.Filename)
		c.SaveUploadedFile(file, "./files"+file.Filename)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
		})
	})
	// if err := seeders.Seed(
	// 	database.DB,
	// 	client.RedisClient,
	// ); err != nil {
	// 	log.Fatal(err)
	// }
	// Configure CORS middleware
	// "github.com/gin-contrib/cors"
	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://localhost:3000"} // Replace with your frontend origin(s)
	// // You can also use AllowAllOrigins: true for development, but be cautious in production
	// // config.AllowAllOrigins = true

	// config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	// config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"} // Add any custom headers your frontend sends
	// config.ExposeHeaders = []string{"Content-Length"} // Headers that the browser can access
	// config.AllowCredentials = true // If your frontend sends cookies or auth headers

	// router.Use(cors.New(config))
	router.Use(middlewares.CORSMiddleware())
	// Pass your GORM db instance to the middleware
	//router.Use(middlewares.AuthenticateAppBySecretKey(database.DB))
	//router.SetTrustedProxies([]string{"192.168.1.2"})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	router.POST(
		"/api/v1/allow",
		middlewares.AuthenticateAppBySecretKey(database.DB),
		controllers.InitiateKyc,
	)
	router.GET("/api/v1/allow/:kyc_token", controllers.FetchKycRequest)

	routes.RegisterAuthRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
