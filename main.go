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
	var x int = 10
	var y byte = 100
	var sum3 int = int(y) + x
	var sum4 byte = byte(x) + y
	fmt.Println(sum3, sum4)

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
	// Pass your GORM db instance to the middleware
	//router.Use(middlewares.AuthenticateAppBySecretKey(database.DB))
	//router.SetTrustedProxies([]string{"192.168.1.2"})

	router.POST(
		"/api/v1/allow",
		middlewares.AuthenticateAppBySecretKey(database.DB),
		controllers.InitiateKyc,
	)
	router.GET("/api/v1/allow/:kyc_token", controllers.FetchKycRequest)
	router.GET("/ping", func(c *gin.Context) {
		day := 4
		switch day {
		case 1:
			fmt.Println("monday")
		case 2:
			fmt.Println("tuesday")
		default:
			fmt.Println("Not a weekday")
		}
		for i := 1; i < 5; i++ {
			fmt.Println(i)
		}
		fruits := [3]string{"apple", "orange", "banana"}
		for _, value := range fruits {
			fmt.Println(value)
		}
		for id, _ := range fruits {
			fmt.Println(id)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

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
