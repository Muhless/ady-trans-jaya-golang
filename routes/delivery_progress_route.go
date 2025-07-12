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
			destinationGroup.POST("/", destinationController.CreateDestinations)
			destinationGroup.GET("/:id", destinationController.GetDestinationsByDeliveryID)
			destinationGroup.POST("", destinationController.UploadArrivalPhoto)
		}
	}
}
