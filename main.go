package main

import (
	"fmt"
	"net/http"

	database "confam-api/database"
	routes "confam-api/routes"
	client "confam-api/utils"
	utils "confam-api/utils"

	middlewares "confam-api/middlewares"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := database.Config{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
	}
	if err := database.Connect(cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database connected.")
	// Call the setup function to register custom validators
	utils.SetupValidator()
	client.InitRedisClient() // Initialize Redis
	router := gin.Default()
	// Set Gin to release mode to reduce noise
	//gin.SetMode(gin.ReleaseMode)

	router.MaxMultipartMemory = 8 << 20 //8Mb
	router.POST("/api/v1/upload", func(c *gin.Context) {
		//form, _ := c.MultipartForm
		//files := form.File["files"]
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

	// Apply middleware
	router.Use(middlewares.CorsMiddleware())
	router.Use(middlewares.RequestLoggingMiddleware())
	router.Use(middlewares.CorsHeadersLoggingMiddleware())
	// Pass your GORM db instance to the middleware
	//router.Use(middlewares.AuthenticateAppBySecretKey(database.DB))
	//router.SetTrustedProxies([]string{"192.168.1.2"})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	routes.RegisterAuthRoutes(router)
	routes.RegisterKycRoutes(router)

	// router.POST(
	// 	"/api/v1/allow",
	// 	middlewares.AuthenticateAppBySecretKey(database.DB),
	// 	controllers.InitiateKyc,
	// )
	// router.GET("/api/v1/allow/:kyc_token", controllers.FetchKycRequest)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
