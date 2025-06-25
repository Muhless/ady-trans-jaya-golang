package main

import (
	"ady-trans-jaya-golang/controllers"
	"ady-trans-jaya-golang/db"
	"ady-trans-jaya-golang/routes"
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
			"http://10.0.2.2",       // android
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
	r.POST("/api/login-driver", controllers.LoginDriver)

	controllers.UserControllers(r, database)
	controllers.DriversControllers(r, database)
	controllers.VehicleControllers(r, database)
	controllers.CustomersControllers(r, database)
	controllers.DeliveryControllers(r, database)

	routes.SetupTransactionRoutes(r, database)
	routes.RegisterDeliveryRoutes(r, database)
	routes.RegisterDeliveryItemRoutes(r, database)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server run error:", err)
	}
}
