package controllers

import (
	"ady-trans-jaya-golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
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

	r.GET("/api/users/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var user model.User
		if err := db.First(&user, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": user})
	})

	r.POST("/api/users", func(ctx *gin.Context) {
		var user model.User

		if err := ctx.ShouldBind(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed save user data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": user})
	})

	r.PUT("/api/users/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var user model.User

		if err := db.First(&user, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": user})
	})
}
