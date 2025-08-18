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

func (c *DeliveryDestinationControllers) UploadPickupPhoto(ctx *gin.Context) {
	c.uploadDeliveryDestinationPhoto(ctx, "pickup_photo_url", "pickup_time")
}

func (c *DeliveryDestinationControllers) UploadDeliveryPhoto(ctx *gin.Context) {
	c.uploadDeliveryDestinationPhoto(ctx, "arrival_photo_url", "arrival_time")
}

type CreateDeliveryDestinationRequest struct {
	DeliveryID        int        `json:"delivery_id" binding:"required"`
	DeliveryStartTime *time.Time `json:"delivery_start_time" binding:"required"`
}

func (c *DeliveryDestinationControllers) GetDestinationByDeliveryID(ctx *gin.Context) {
	deliveryID := ctx.Param("id")
	var progresses []model.DeliveryDestination

	if err := c.DB.Where("delivery_id = ?", deliveryID).Find(&progresses).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data progress",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    progresses,
		"message": "Berhasil mengambil data progress",
	})
}

func (c *DeliveryDestinationControllers) CreateDeliveryDestination(ctx *gin.Context) {
	var req CreateDeliveryDestinationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing model.DeliveryDestination
	if err := db.DB.Where("delivery_id = ?", req.DeliveryID).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Progress untuk pengiriman ini sudah dibuat"})
		return
	}

	progress := model.DeliveryDestination{
		DeliveryID:        req.DeliveryID,
		DeliveryStartTime: req.DeliveryStartTime,
	}

	if err := db.DB.Create(&progress).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan progress"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Progress pengiriman berhasil dibuat"})
}

func (c *DeliveryDestinationControllers) uploadDeliveryDestinationPhoto(ctx *gin.Context, field string, timeField string) {
	id := ctx.Param("id")

	var progress model.DeliveryDestination
	if err := db.DB.First(&progress, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
		return
	}

	file, err := ctx.FormFile("photo")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	filename := fmt.Sprintf("%s_%d_%d.jpg", field, progress.ID, time.Now().Unix())
	path := fmt.Sprintf("uploads/%s", filename)
	if err := ctx.SaveUploadedFile(file, path); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		return
	}

	url := fmt.Sprintf("/%s", path)

	updates := map[string]interface{}{
		field:        url,
		timeField:    time.Now(),
		"updated_at": time.Now(),
	}

	if err := db.DB.Model(&progress).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Upload success", "url": url})
	fmt.Println("File berhasil diupload:", filename)
}

func (c *DeliveryDestinationControllers) DeleteDeliveryDestination(ctx *gin.Context) {
	id := ctx.Param("id")

	var progress model.DeliveryDestination
	if err := db.DB.First(&progress, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Progress tidak ditemukan"})
		return
	}

	if err := db.DB.Delete(&progress).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus progress"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Progress berhasil dihapus"})
}

func (c *DeliveryDestinationControllers) UploadArrivalPhoto(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid destination ID"})
		return
	}

	var destination model.DeliveryDestination
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
