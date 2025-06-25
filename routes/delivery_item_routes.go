package routes

import (
	"ady-trans-jaya-golang/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterDeliveryItemRoutes(router *gin.Engine, db *gorm.DB) {
	controller := controllers.DeliveryItemController{DB: db}

	api := router.Group("/api")
	{
		itemGroup := api.Group("/delivery-items")
		{
			itemGroup.POST("", controller.CreateDeliveryItem)
			itemGroup.GET("", controller.GetDeliveryItems)
			itemGroup.GET("/:id", controller.GetDeliveryItem)
			itemGroup.PUT("/:id", controller.UpdateDeliveryItem)
			itemGroup.DELETE("/:id", controller.DeleteDeliveryItem)
		}
	}
}
