package controllers

import (
	"ady-trans-jaya-golang/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionController struct {
	DB *gorm.DB
}

func NewTransactionController(db *gorm.DB) *TransactionController {
	return &TransactionController{DB: db}
}

// GetTransactions - Get all transactions
func (c *TransactionController) GetTransactions(ctx *gin.Context) {
	var transactions []model.Transaction

	if err := c.DB.
		Preload("Customer").
		Preload("Delivery").
		Find(&transactions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data transaksi",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}

// GetTransactionByID - Get transaction by ID
func (c *TransactionController) GetTransactionByID(ctx *gin.Context) {
	id := ctx.Param("id")
	
	var transaction model.Transaction
	if err := c.DB.
		Preload("Customer").
		Preload("Delivery").
		First(&transaction, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Transaksi tidak ditemukan",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data transaksi",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": transaction,
	})
}

// CreateTransaction - Create new transaction
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
			fmt.Println("Recovered in CreateTransaction:", r)
		}
	}()

	request.Transaction.TotalDelivery = len(request.Deliveries)

	if err := tx.Create(&request.Transaction).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save transaction data"})
		return
	}

	for _, delivery := range request.Deliveries {
		delivery.TransactionID = request.Transaction.ID
		if err := tx.Omit("Driver", "Vehicle").Create(&delivery).Error; err != nil {
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
	if err := c.DB.Preload("Customer").Preload("Delivery").First(&result, request.Transaction.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created transaction with related data"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Transaction data successfully saved",
		"data":    result,
	})
}

// UpdateTransaction - Update existing transaction
func (c *TransactionController) UpdateTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	
	type UpdateTransactionRequest struct {
		model.Transaction
		Deliveries []model.Delivery `json:"deliveries"`
	}

	var request UpdateTransactionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if transaction exists
	var existingTransaction model.Transaction
	if err := c.DB.First(&existingTransaction, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Transaksi tidak ditemukan",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data transaksi",
		})
		return
	}

	tx := c.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in UpdateTransaction:", r)
		}
	}()

	// Update transaction basic info
	request.Transaction.ID = existingTransaction.ID
	request.Transaction.TotalDelivery = len(request.Deliveries)
	
	if err := tx.Model(&existingTransaction).Updates(request.Transaction).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction data"})
		return
	}

	// Delete existing deliveries
	if err := tx.Where("transaction_id = ?", existingTransaction.ID).Delete(&model.Delivery{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old delivery data"})
		return
	}

	// Create new deliveries
	for _, delivery := range request.Deliveries {
		delivery.ID = 0 // Reset ID for new records
		delivery.TransactionID = existingTransaction.ID
		if err := tx.Omit("Driver", "Vehicle").Create(&delivery).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save delivery data"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Return updated transaction with relations
	var result model.Transaction
	if err := c.DB.Preload("Customer").Preload("Delivery").First(&result, existingTransaction.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated transaction with related data"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transaction data successfully updated",
		"data":    result,
	})
}

// DeleteTransaction - Delete transaction by ID
func (c *TransactionController) DeleteTransaction(ctx *gin.Context) {
	id := ctx.Param("id")

	// Check if transaction exists
	var transaction model.Transaction
	if err := c.DB.First(&transaction, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Transaksi tidak ditemukan",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data transaksi",
		})
		return
	}

	tx := c.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in DeleteTransaction:", r)
		}
	}()

	// Delete related deliveries first
	if err := tx.Where("transaction_id = ?", id).Delete(&model.Delivery{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete delivery data"})
		return
	}

	// Delete transaction
	if err := tx.Delete(&transaction).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction data"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transaction data successfully deleted",
	})
}

// GetTransactionsByCustomer - Get transactions by customer ID
func (c *TransactionController) GetTransactionsByCustomer(ctx *gin.Context) {
	customerID := ctx.Param("customer_id")
	
	var transactions []model.Transaction
	if err := c.DB.
		Preload("Customer").
		Preload("Delivery").
		Where("customer_id = ?", customerID).
		Find(&transactions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data transaksi",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}

// GetTransactionsPaginated - Get transactions with pagination
func (c *TransactionController) GetTransactionsPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	
	offset := (page - 1) * limit
	
	var transactions []model.Transaction
	var total int64
	
	// Count total records
	if err := c.DB.Model(&model.Transaction{}).Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghitung total data",
		})
		return
	}
	
	// Get paginated data
	if err := c.DB.
		Preload("Customer").
		Preload("Delivery").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data transaksi",
		})
		return
	}
	
	totalPages := (int(total) + limit - 1) / limit
	
	ctx.JSON(http.StatusOK, gin.H{
		"data": transactions,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// SearchTransactions - Search transactions by various criteria
func (c *TransactionController) SearchTransactions(ctx *gin.Context) {
	search := ctx.Query("search")
	status := ctx.Query("status")
	
	query := c.DB.
		Preload("Customer").
		Preload("Delivery")
	
	if search != "" {
		query = query.Joins("LEFT JOIN customers ON transactions.customer_id = customers.id").
			Where("customers.name LIKE ? OR transactions.id LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	
	if status != "" {
		query = query.Where("transactions.status = ?", status)
	}
	
	var transactions []model.Transaction
	if err := query.Find(&transactions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mencari data transaksi",
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}

// UpdateTransactionStatus - Update only transaction status
func (c *TransactionController) UpdateTransactionStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	
	type UpdateStatusRequest struct {
		Status string `json:"status" binding:"required"`
	}
	
	var request UpdateStatusRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var transaction model.Transaction
	if err := c.DB.First(&transaction, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Transaksi tidak ditemukan",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data transaksi",
		})
		return
	}
	
	if err := c.DB.Model(&transaction).Update("status", request.Status).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengupdate status transaksi",
		})
		return
	}
	
	// Return updated transaction
	if err := c.DB.Preload("Customer").Preload("Delivery").First(&transaction, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data transaksi yang sudah diupdate",
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Status transaksi berhasil diupdate",
		"data":    transaction,
	})
}