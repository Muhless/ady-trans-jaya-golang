package controllers

import (
	"ady-trans-jaya-golang/model"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func VehicleControllers(r *gin.Engine, db *gorm.DB) {
	// Get all vehicles
	r.GET("/api/vehicle", func(ctx *gin.Context) {
		var vehicle []model.Vehicle
		if err := db.Find(&vehicle).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vehicle data"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": vehicle})
	})

	// Get vehicle by ID
	r.GET("/api/vehicle/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var vehicle model.Vehicle
		if err := db.First(&vehicle, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": vehicle})
	})

	// Create a new vehicle
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

		ctx.JSON(http.StatusOK, gin.H{"message": "Vehicle Data Successfully Saved", "data": vehicle})
	})

	// Update vehicle by ID
	r.PUT("/api/vehicle/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var vehicle model.Vehicle

		// Get the existing vehicle by ID
		if err := db.First(&vehicle, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			return
		}

		// Bind the incoming data to the vehicle struct
		if err := ctx.ShouldBindJSON(&vehicle); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save the updated vehicle
		if err := db.Save(&vehicle).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Vehicle Data Successfully Updated", "data": vehicle})
	})

	// Delete vehicle by ID
	r.DELETE("/api/vehicle/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var vehicle model.Vehicle

		// Find the vehicle by ID
		if err := db.First(&vehicle, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			return
		}

		// Delete the vehicle
		if err := db.Delete(&vehicle).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vehicle"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Vehicle successfully deleted"})
	})

	// Search vehicles (by name or license plate)
	r.GET("/api/vehicle/search", func(ctx *gin.Context) {
		searchQuery := ctx.DefaultQuery("query", "")
		if searchQuery == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
			return
		}

		var vehicles []model.Vehicle
		// Searching by name or license plate (case insensitive)
		if err := db.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(searchQuery)+"%").
			Or("LOWER(license_plate) LIKE ?", "%"+strings.ToLower(searchQuery)+"%").
			Find(&vehicles).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search vehicles"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": vehicles})
	})
}
