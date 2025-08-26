package main

import (
	"fmt"
	"net/http"
	"sync"

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

var once sync.Once

type singleton struct {
	Data string
}

var instance *singleton

func GetInstance() *singleton {
	once.Do(func() {
		instance = &singleton{Data: "This is a string"}
	})
	return instance
}

func addTo(base int, vals ...int) []int {
	out := make([]int, 0, len(vals))
	for _, v := range vals {
		out = append(out, base+v)
	}
	return out
}

func nonempty(strings []string) []string {
	i := 0
	for _, s := range strings {
		if s != "" {
			strings[i] = s
			i++
		}
	}
	return strings[:i]
}

func main() {
	godotenv.Load()
	// Call the setup function to register custom validators
	utils.SetupValidator()
	client.InitRedisClient() // Initialize Redis
	router := gin.Default()
	// Set Gin to release mode to reduce noise
	//gin.SetMode(gin.ReleaseMode)
	database.ConnectDatabase()

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

	s1 := GetInstance()
	s2 := GetInstance()

	fmt.Printf("s1: %p\n", s1) // %p prints the memory address
	fmt.Printf("s2: %p\n", s2)
	fmt.Printf("Are they the same instance? %t\n", s1 == s2)
	fmt.Println(s1.Data)

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
