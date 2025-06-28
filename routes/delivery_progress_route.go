package routes

import (
	"ady-trans-jaya-golang/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterDeliveryProgressRoutes(router *gin.Engine, db *gorm.DB) {
	progressController := controllers.DeliveryProgressController{DB: db}

	api := router.Group("/api")
	{
		progressGroup := api.Group("/delivery-progress")
		{
			progressGroup.GET("/:id/progress", progressController.GetProgressByDeliveryID)
			progressGroup.POST("", progressController.CreateDeliveryProgress)
			progressGroup.DELETE("/:id", progressController.DeleteDeliveryProgress)

			// Tambahan untuk upload foto
			progressGroup.POST("/upload-pickup/:id", progressController.UploadPickupPhoto)
			progressGroup.POST("/upload-delivery/:id", progressController.UploadDeliveryPhoto)
		}
	}
}
