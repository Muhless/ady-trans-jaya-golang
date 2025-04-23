package controllers

import (
	"ady-trans-jaya-golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TransactionController(r *gin.Engine, db *gorm.DB) {
	r.GET("/api/transaction", func(ctx *gin.Context) {
		var transaction []model.Transaction
		if err := db.Find(&transaction).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed tp retrieve transactions data"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": transaction})
	})

	r.POST("/api/transaction", func(ctx *gin.Context) {
		var transaction model.Transaction

		if err := ctx.ShouldBindJSON(&transaction); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&transaction).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save transaction data"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Transaction data successfully saved", "data": transaction})
	})

}
