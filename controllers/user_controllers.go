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
	"gorm.io/gorm"
)

func UserControllers(r *gin.Engine, db *gorm.DB) {
	r.GET("/api/users", func(ctx *gin.Context) {
		var user []model.User
		if err := db.Find(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users data"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": user})
	})

	r.GET("/api/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var user model.User
		if err := db.First(&user, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": user})
	})

	r.POST("/api/user", func(ctx *gin.Context) {
		var user model.User

		if err := ctx.ShouldBind(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)

		if err := db.Create(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed save user data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": user})
	})

	r.PUT("/api/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var existingUser model.User

		if err := db.First(&existingUser, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		var input model.User
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		existingUser.Username = input.Username
		existingUser.Role = input.Role

		if input.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
				return
			}
			existingUser.Password = string(hashedPassword)
		}

		if err := db.Save(&existingUser).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user data"})
			return
		}

		existingUser.Password = ""
		ctx.JSON(http.StatusOK, gin.H{"data": existingUser})
	})

}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Input tidak valid"})
		return
	}

	var user model.User
	if err := db.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Username tidak ditemukan"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Password salah"})
		return
	}

	// Buat JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(2 * time.Hour).Unix(), // expired 2 jam
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
		"token": tokenString,
	})
}
