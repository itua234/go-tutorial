package main

import (
	"fmt"
	"net/http"

	"blog/config"

	"github.com/gin-gonic/gin"
)

func main() {
	var firstname, lastname string = "itua", "osemeilu"
	var age int = 18
	length := 20
	var (
		a int
		b        = 1
		c string = "hello"
	)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(firstname + lastname)
	fmt.Println("Length is:", length)
	fmt.Println(age)

	router := gin.Default()
	config.ConnectDatabase()
	//router.SetTrustedProxies([]string{"192.168.1.2"})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.Run(":8000") //listen and serve on 0.0.0.0:8000
}
