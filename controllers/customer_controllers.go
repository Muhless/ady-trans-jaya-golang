package controllers

import (
	"ady-trans-jaya-golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CustomersControllers(r *gin.Engine, db *gorm.DB) {
	r.GET("/api/customers", func(ctx *gin.Context) {
		var customer []model.Customer
		if err := db.Find(&customer).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers data"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": customer})
	})

	r.GET("/api/customer/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var customer model.Customer
		if err := db.First(&customer, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": customer})
	})

	r.POST("/api/customer", func(ctx *gin.Context) {
		var customer model.Customer

		if err := ctx.ShouldBindJSON(&customer); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&customer).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save customer data"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Customer data successfully saved", "data": customer})
	})

}
