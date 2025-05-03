package controllers

import (
	"ady-trans-jaya-golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionController struct {
	DB *gorm.DB
}

func NewTransactionController(db *gorm.DB) *TransactionController {
	return &TransactionController{DB: db}
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	type CreateTransactionRequest struct {
		model.Transaction
		Deliveries []model.Delivery `json:"deliveries"`
	}

	var request CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := c.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	request.Transaction.TotalDelivery = len(request.Deliveries)

	if err := tx.Create(&request.Transaction).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save transaction data"})
		return
	}

	for i := range request.Deliveries {
		// Validasi wajib isi dan format delivery_date
		if request.Deliveries[i].DeliveryDate.IsZero() {
			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "delivery_date is required"})
			return
		}

		request.Deliveries[i].TransactionID = request.Transaction.ID
		if err := tx.Create(&request.Deliveries[i]).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save delivery data"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	var result model.Transaction
	if err := c.DB.Preload("Customer").Preload("Deliveries").First(&result, request.Transaction.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created transaction with related data"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transaction data successfully saved",
		"data":    result,
	})
}
