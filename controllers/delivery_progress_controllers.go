package controllers

import (
	"ady-trans-jaya-golang/db"
	"ady-trans-jaya-golang/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeliveryProgressController struct {
	DB *gorm.DB
}

func (c *DeliveryProgressController) UploadPickupPhoto(ctx *gin.Context) {
	c.uploadDeliveryProgressPhoto(ctx, "pickup_photo_url", "pickup_time")
}

func (c *DeliveryProgressController) UploadDeliveryPhoto(ctx *gin.Context) {
	c.uploadDeliveryProgressPhoto(ctx, "arrival_photo_url", "arrival_time")
}

func (c *DeliveryProgressController) GetProgressByDeliveryID(ctx *gin.Context) {
	deliveryID := ctx.Param("id")
	var progresses []model.DeliveryProgress

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

func (c *DeliveryProgressController) CreateDeliveryProgress(ctx *gin.Context) {
	var req model.DeliveryProgress
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing model.DeliveryProgress
	if err := db.DB.Where("delivery_id = ?", req.DeliveryID).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "Progress untuk pengiriman ini sudah dibuat",
		})
		return
	}

	progress := model.DeliveryProgress{
		DeliveryID:        req.DeliveryID,
		DeliveryStartTime: req.DeliveryStartTime,
	}

	if err := db.DB.Create(&progress).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan progress"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Progress pengiriman berhasil dibuat"})
}

func (c *DeliveryProgressController) uploadDeliveryProgressPhoto(ctx *gin.Context, field string, timeField string) {
	id := ctx.Param("id")

	var progress model.DeliveryProgress
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
}

func (c *DeliveryProgressController) DeleteDeliveryProgress(ctx *gin.Context) {
	id := ctx.Param("id")

	var progress model.DeliveryProgress
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
