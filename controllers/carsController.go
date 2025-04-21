package controllers

import (
	"ady-trans-jaya-golang/model"
	"bytes"
	"fmt"
	"io"
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

		body, _ := io.ReadAll(ctx.Request.Body)
		fmt.Println("📥 Raw Body:", string(body))
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		if err := ctx.ShouldBindJSON(&cars); err != nil {
			fmt.Println("❌ Error binding:", err.Error())
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
