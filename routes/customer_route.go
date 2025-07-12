package routes

import (
	"ady-trans-jaya-golang/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterCustomerRoutes(router *gin.Engine, db *gorm.DB) {
	controller := controllers.CustomersController{DB: db}

	api := router.Group("/api")
	{
		customerGroup := api.Group("/customer")
		{
			customerGroup.PATCH("/:id", controller.UpdateCustomer)
		}
	}
}
