package controllers

import (
	"ady-trans-jaya-golang/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeliveryProgressController struct {
	DB *gorm.DB
}

func NewDeliveryProgressController(db *gorm.DB) *DeliveryProgressController {
	return &DeliveryProgressController{DB: db}
}

func (c *DeliveryProgressController) GetProgressByDeliveryID(ctx *gin.Context) {
	id := ctx.Param("id")
	var progress []model.DeliveryProgress

	if err := c.DB.Where("delivery_id = ?", id).Find(&progress).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data progres"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": progress})
}

func (c *DeliveryProgressController) CreateProgress(ctx *gin.Context) {
	deliveryIDStr := ctx.Param("id")
	deliveryID, err := strconv.Atoi(deliveryIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var input model.DeliveryProgress
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	input.DeliveryID = deliveryID
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	if err := c.DB.Create(&input).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan progres pengiriman"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Progres pengiriman berhasil ditambahkan", "data": input})
}

func (c *DeliveryProgressController) PatchProgress(ctx *gin.Context) {
	id := ctx.Param("id")
	var updateFields map[string]interface{}

	if err := ctx.ShouldBindJSON(&updateFields); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	if err := c.DB.Model(&model.DeliveryProgress{}).
		Where("id = ?", id).
		Updates(updateFields).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate progress"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Progress berhasil diperbarui"})
}
