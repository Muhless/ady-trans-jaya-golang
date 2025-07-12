package controllers

import (
	"ady-trans-jaya-golang/db"
	"ady-trans-jaya-golang/model"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeliveryDestinationControllers struct {
	DB *gorm.DB
}

func (c *DeliveryDestinationControllers) CreateDestinations(ctx *gin.Context) {
	var destinations []model.DeliveryDestinations

	if err := ctx.ShouldBindJSON(&destinations); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	for i, dest := range destinations {
		if dest.DeliveryID == 0 || dest.Address == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid destination at index %d: delivery_id and address required", i),
			})
			return
		}
	}

	// Simpan semua ke database
	if err := db.DB.Create(&destinations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create destinations"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":      "Destinations created successfully",
		"destinations": destinations,
	})
}

func (c *DeliveryDestinationControllers) GetDestinationsByDeliveryID(ctx *gin.Context) {
	deliveryIDStr := ctx.Param("id")
	deliveryID, err := strconv.Atoi(deliveryIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid delivery ID"})
		return
	}

	var destinations []model.DeliveryDestinations
	if err := c.DB.
		Where("delivery_id = ?", deliveryID).
		Order("sequence ASC").
		Find(&destinations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get destinations"})
		return
	}

	ctx.JSON(http.StatusOK, destinations)
}

func (c *DeliveryDestinationControllers) UploadArrivalPhoto(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid destination ID"})
		return
	}

	var destination model.DeliveryDestinations
	if err := db.DB.First(&destination, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Destination not found"})
		return
	}

	file, err := ctx.FormFile("photo")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("arrival_destination_%d_%d.jpg", destination.ID, timestamp)
	path := fmt.Sprintf("uploads/%s", filename)
	url := fmt.Sprintf("/%s", path)

	if err := ctx.SaveUploadedFile(file, path); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	now := time.Now()
	update := map[string]interface{}{
		"arrival_photo_url": url,
		"arrival_time":      now,
		"updated_at":        now,
		"status":            "selesai",
	}
	if err := db.DB.Model(&destination).Updates(update).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update destination"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Photo uploaded successfully",
		"url":     url,
	})
}
