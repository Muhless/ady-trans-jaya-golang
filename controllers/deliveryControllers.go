package controllers

import (
	"ady-trans-jaya-golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DeliveryControllers(r *gin.Engine, db *gorm.DB) {
	r.GET("/api/delivery", func(ctx *gin.Context) {
		var delivery []model.Delivery
		if err := db.Find(&delivery).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve delivery data"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": delivery})
	})

	r.POST("/api/delivery", func(ctx *gin.Context) {
		var delivery model.Delivery
		if err := ctx.ShouldBindJSON(&delivery); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if delivery.DeliveryStatus == "" {
			delivery.Driver.Status = "menunggu persetujuan"
		}

		if err := db.Create(&delivery).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save delivery data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Delivery data successfully saved", "data": delivery})
	})

}
