package main

import (
	"ady-trans-jaya-golang/controllers"
	"ady-trans-jaya-golang/db"
	"log"
	"time"

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
		AllowOrigins: []string{
			"http://localhost:5173", // untuk development
			"http://202.10.41.13",   // alamat VPS frontend
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes
	r.POST("/api/login", controllers.Login)
	controllers.UserControllers(r, database)
	controllers.DriversControllers(r, database)
	controllers.VehicleControllers(r, database)
	controllers.CustomersControllers(r, database)

	transactionController := controllers.NewTransactionController(database)
	r.POST("/api/transactions", transactionController.CreateTransaction)
	r.GET("/api/transactions", transactionController.GetTransactions)

	controllers.DeliveryControllers(r, database)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server run error:", err)
	}
}
