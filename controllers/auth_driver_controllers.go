package controllers

import (
	"ady-trans-jaya-golang/db"
	"ady-trans-jaya-golang/model"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func LoginDriver(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Input tidak valid"})
		return
	}

	var driver model.Driver
	if err := db.DB.Where("username = ?", input.Username).First(&driver).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Username tidak ditemukan"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(driver.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Password salah"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       driver.ID,
		"username": driver.Username,
		"exp":      time.Now().Add(2 * time.Hour).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"driver": gin.H{
			"id":       driver.ID,
			"name":     driver.Name,
			"username": driver.Username,
			"phone":    driver.Phone,
			"address":  driver.Address,
			"photo":    driver.Photo,
		},
	})
}
