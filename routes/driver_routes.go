package routes

import (
	"ady-trans-jaya-golang/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterDriverRoutes(router *gin.Engine, db *gorm.DB) {
	controller := controllers.DriverController{DB: db}
	api := router.Group("/api")
	{
		driverGroup := api.Group("/driver")
		{
			driverGroup.PATCH("/:id/active", controller.GetActiveDeliveryForDriver)

		}
	}
}
