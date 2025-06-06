package controllers

import (
	"ady-trans-jaya-golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserControllers(r *gin.Engine, db *gorm.DB) {
	r.GET("/api/users", func(ctx *gin.Context) {
		var user []model.User
		db.Order("id ASC").Find(&user)

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
