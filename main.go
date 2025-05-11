package main

import (
	"ady-trans-jaya-golang/controllers"
	"ady-trans-jaya-golang/db"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Authorization", "Content-Type"},
	}))

	controllers.UserControllers(r, database)
	controllers.DriversControllers(r, database)
	controllers.VehicleControllers(r, database)
	controllers.CustomersControllers(r, database)
	transactionController := controllers.NewTransactionController(database)
	r.POST("/api/transactions", transactionController.CreateTransaction)
	r.GET("/api/transactions", transactionController.GetTransactions)
	// r.POST("/api/users", handler.LoginHandler(database))
	controllers.DeliveryControllers(r, database)

	r.Run(":8080")
}
