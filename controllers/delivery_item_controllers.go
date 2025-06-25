package controllers

import (
	"ady-trans-jaya-golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeliveryItemController struct {
	DB *gorm.DB
}

func (c *DeliveryItemController) CreateDeliveryItem(ctx *gin.Context) {
	var item model.DeliveryItem
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.DB.Create(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}
	ctx.JSON(http.StatusCreated, item)
}

func (c *DeliveryItemController) GetDeliveryItems(ctx *gin.Context) {
	var items []model.DeliveryItem
	if err := c.DB.Find(&items).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func (c *DeliveryItemController) GetDeliveryItem(ctx *gin.Context) {
	id := ctx.Param("id")
	var item model.DeliveryItem
	if err := c.DB.First(&item, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *DeliveryItemController) UpdateDeliveryItem(ctx *gin.Context) {
	id := ctx.Param("id")
	var item model.DeliveryItem

	if err := c.DB.First(&item, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.DB.Save(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *DeliveryItemController) DeleteDeliveryItem(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&model.DeliveryItem{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
