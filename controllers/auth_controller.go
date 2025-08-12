package controllers

import (
	"net/http"
	"os"
	"time"

	database "confam-api/database"
	models "confam-api/models"
	structs "confam-api/structs"
	utils "confam-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var req structs.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Registration Successful",
	})
}

func Login(c *gin.Context) {
	var req structs.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, structs.ValidationError{
			Status:  "failed",
			Message: "Validation failed",
			//Errors:  formatBindingError(err),
		})
		return
	}

	// Find the company by email
	var company models.Company
	result := database.DB.Where("email = ?", req.Email).First(&company)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error."})
		return
	}

	// Compare the password with the stored hash
	if err := utils.ComparePasswordAndHash(req.Password, *company.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials."})
		return
	}

	// Generate a JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": company.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login Successful",
		"data":    tokenString,
	})
}
