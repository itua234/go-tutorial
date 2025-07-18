package main

import (
	"fmt"
	"net/http"

	"confam-api/database"

	"confam-api/middlewares"
	client "confam-api/utils"

	"confam-api/controllers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Login struct {
	User     string
	Password string
}

func main() {
	x := []string{"a", "b", "c", "d"}
	y := x[:2]
	z := x[1:]
	x[1] = "y"
	y[0] = "x"
	z[1] = "z"
	fmt.Println("x:", x)
	fmt.Println("y:", y)
	fmt.Println("z:", z)

	godotenv.Load()
	client.InitRedisClient() // Initialize Redis
	router := gin.Default()
	database.ConnectDatabase()
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

	router.POST(
		"/api/v1/allow",
		middlewares.AuthenticateAppBySecretKey(database.DB),
		controllers.InitiateKyc,
	)
	router.GET("/api/v1/allow/:kyc_token", controllers.FetchKycRequest)

	router.POST("/hello", func(c *gin.Context) {
		//name := c.PostForm("name") // for form data
		// name := c.Query("name") // for query string
		var json Login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Hello " + json.Password})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
