package routes

import (
	"ady-trans-jaya-golang/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterDeliveryRoutes(router *gin.Engine, db *gorm.DB) {
	controller := controllers.DeliveryController{DB: db}

	api := router.Group("/api")
	{
		deliveryGroup := api.Group("/delivery")
		{
			deliveryGroup.PATCH("/:id/status", controller.UpdateDeliveryStatus)
			deliveryGroup.GET("/driver/:id/active", controllers.GetActiveDeliveriesByDriver)
			deliveryGroup.PATCH("/driver/:id/status", controller.UpdateDeliveryDriverStatus)
		}
	}
}
