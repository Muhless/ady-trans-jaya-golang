package controllers

import (
	"ady-trans-jaya-golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterCustomerRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/api/customer", func(ctx *gin.Context) {
		var customer []model.Customer
		if err := db.Find(&customer).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers data"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": customer})
	})

}
