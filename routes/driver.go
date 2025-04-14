package routes

import (
	"ady-trans-jaya-golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterDriverRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/api/driver", func(ctx *gin.Context) {
		var driver []model.Driver
		if err := db.Find(&driver).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve drivers data"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": driver})
	})

	router.POST("/api/driver", func(ctx *gin.Context) {
		var driver model.Driver
		if err := ctx.ShouldBindBodyWithJSON(&driver); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&driver).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed save drivers data"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Driver successfully saved"})
	})
}
