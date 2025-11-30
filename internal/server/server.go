package server

import (
	"confam-api/internal/middlewares"
	client "confam-api/internal/redis"
	routes "confam-api/internal/routes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewRouter initializes the Gin engine with all middleware and routes
func NewRouter(rdb *client.Client) *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8Mb
	router.Use(middlewares.CorsMiddleware())
	router.Use(middlewares.RequestLoggingMiddleware())
	router.Use(middlewares.CorsHeadersLoggingMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	routes.RegisterAuthRoutes(router)
	routes.RegisterKycRoutes(router, rdb)
	router.POST("/api/v1/upload", handleUpload)
	//router.Use(middlewares.AuthenticateAppBySecretKey(database.DB))
	//router.SetTrustedProxies([]string{"192.168.1.2"})

	// router.POST(
	// 	"/api/v1/allow",
	// 	middlewares.AuthenticateAppBySecretKey(database.DB),
	// 	controllers.InitiateKyc,
	// )
	// router.GET("/api/v1/allow/:kyc_token", controllers.FetchKycRequest)

	return router
}

func handleUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	// Ensure the directory exists
	// os.MkdirAll("./files", os.ModePerm)

	dst := "./files/" + file.Filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
	})
}
