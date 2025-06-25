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
			progressGroup.GET("/:delivery_id", progressController.GetProgressByDeliveryID)
			progressGroup.POST("/", progressController.CreateProgress)
			progressGroup.PATCH("/:id", progressController.PatchProgress)

		}
	}

}
