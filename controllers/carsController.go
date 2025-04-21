package controllers

import (
	"ady-trans-jaya-golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CarsControllers(r *gin.Engine, db *gorm.DB) {
	r.GET("/api/cars", func(ctx *gin.Context) {
		var cars []model.Car
		if err := db.Find(&cars).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cars data"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": cars})
	})

	r.POST("/api/cars", func(ctx *gin.Context) {
		var cars model.Car
		if err := ctx.ShouldBind(&cars); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if cars.Status == "" {
			cars.Status = "tersedia"
		}

		if err := db.Create(&cars).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cars data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Cars Data Successfully Saved", "data": cars})
	})
}
