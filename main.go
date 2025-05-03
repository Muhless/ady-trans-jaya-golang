package main

import (
	"ady-trans-jaya-golang/controllers"
	"ady-trans-jaya-golang/db"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	r := gin.Default()
	r.Use(cors.Default())
	controllers.DriversControllers(r, db)
	controllers.VehicleControllers(r, db)
	controllers.CustomersControllers(r, db)
	transactionController := controllers.NewTransactionController(db)
	r.POST("/api/transactions", transactionController.CreateTransaction)
	r.GET("/api/transactions", transactionController.GetTransactions)
	controllers.DeliveryControllers(r, db)
	r.Run(":8080")
}
