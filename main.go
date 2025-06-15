package main

import (
	"fmt"
	"net/http"

	"blog/config"

	client "blog/util"

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
	fmt.Println(firstname + lastname)
	fmt.Println("Length is:", length)
	fmt.Println(age)

	godotenv.Load()
	client.InitRedisClient() // Initialize Redis
	router := gin.Default()
	config.ConnectDatabase()
	//router.SetTrustedProxies([]string{"192.168.1.2"})

	router.GET("/ping", func(c *gin.Context) {
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
	router.Run(":8000") //listen and serve on 0.0.0.0:8000
}
