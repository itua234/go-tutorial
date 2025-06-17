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
	var firstname, lastname string = "itua", "osemeilu"
	var age int = 18
	length := 20
	var arr1 = [3]int{1, 2, 3}
	//arr2 := [5]int{4,5,6,7,8}
	fmt.Println(arr1)
	fmt.Println(firstname + lastname)
	fmt.Println("Length is:", length)
	fmt.Println(age)

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
	router.Use(middlewares.AuthenticateAppBySecretKey(database.DB))
	//router.SetTrustedProxies([]string{"192.168.1.2"})

	router.POST("/api/v1/allow", controllers.InitiateKyc)
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
	router.GET("/redis-ping", func(c *gin.Context) {
		pong, err := client.RedisClient.Ping(c).Result()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": pong})
	})
	router.POST("/hello", func(c *gin.Context) {
		//name := c.PostForm("name") // for form data
		// name := c.Query("name") // for query string
		//body, _ := c.GetRawData() //no recommended for complex data
		var json Login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Hello " + json.Password})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Server starting on port :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
