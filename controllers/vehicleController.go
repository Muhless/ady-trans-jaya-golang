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

func VehicleControllers(r *gin.Engine, db *gorm.DB) {
	r.GET("/api/vehicle", func(ctx *gin.Context) {
		var vehicle []model.Vehicle
		if err := db.Find(&vehicle).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vehicle data"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": vehicle})
	})

	r.POST("/api/vehicle", func(ctx *gin.Context) {
		var vehicle model.Vehicle

		body, _ := io.ReadAll(ctx.Request.Body)
		fmt.Println("📥 Raw Body:", string(body))
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		if err := ctx.ShouldBindJSON(&vehicle); err != nil {
			fmt.Println("❌ Error binding:", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if vehicle.Status == "" {
			vehicle.Status = "tersedia"
		}

		if err := db.Create(&vehicle).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save vehicle data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "vehicle Data Successfully Saved", "data": vehicle})
	})
}
