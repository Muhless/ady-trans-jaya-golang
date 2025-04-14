package main

import (
	"ady-trans-jaya-golang/internal/db"
	"ady-trans-jaya-golang/routes"
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
	routes.RegisterDriverRoutes(r, db)
	r.Run(":8080")
}
