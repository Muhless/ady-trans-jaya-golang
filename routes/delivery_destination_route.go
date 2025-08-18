package routes

import (
	"ady-trans-jaya-golang/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterDeliveryProgressRoutes(router *gin.Engine, db *gorm.DB) {
	destinationController := controllers.DeliveryDestinationControllers{DB: db}

	api := router.Group("/api")
	{
		destinationGroup := api.Group("/delivery-destination")
		{
			destinationGroup.POST("/", destinationController.CreateDeliveryDestination)
			destinationGroup.GET("/:id/destination", destinationController.GetDestinationByDeliveryID)
			destinationGroup.DELETE("/:id", destinationController.DeleteDeliveryDestination)

			destinationGroup.POST("/upload-pickup/:id", destinationController.UploadPickupPhoto)
			destinationGroup.POST("/upload-delivery/:id", destinationController.UploadDeliveryPhoto)
		}
	}
}
