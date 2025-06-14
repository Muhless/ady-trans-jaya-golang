package routes

import (
	"ady-trans-jaya-golang/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTransactionRoutes(router *gin.Engine, db *gorm.DB) {
	transactionController := controllers.NewTransactionController(db)

	api := router.Group("/api")
	{
		transactions := api.Group("/transactions")
		{
			transactions.GET("", transactionController.GetTransactions)
			transactions.GET("/paginated", transactionController.GetTransactionsPaginated)
			transactions.GET("/search", transactionController.SearchTransactions)
			transactions.GET("/:id", transactionController.GetTransactionByID)
			transactions.POST("", transactionController.CreateTransaction)
			transactions.PUT("/:id", transactionController.UpdateTransaction)
			transactions.DELETE("/:id", transactionController.DeleteTransaction)
			transactions.PATCH("/:id/status", transactionController.UpdateTransactionStatus)
		}

		customers := api.Group("/customers")
		{
			customers.GET("/:customer_id/transactions", transactionController.GetTransactionsByCustomer)
		}
	}
}
