package controllers

import (
	"ady-trans-jaya-golang/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionController struct {
	DB *gorm.DB
}

func NewTransactionController(db *gorm.DB) *TransactionController {
	return &TransactionController{DB: db}
}

func (c *TransactionController) GetTransactions(ctx *gin.Context) {
	var transactions []model.Transaction

	if err := c.DB.
		Preload("Customer").
		Preload("Delivery").
		Preload("Delivery.Driver").
		Preload("Delivery.Vehicle").
		Preload("Delivery.Items").
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

	request.Transaction.ID = existingTransaction.ID
	request.Transaction.TotalDelivery = len(request.Deliveries)

	if err := tx.Model(&existingTransaction).Updates(request.Transaction).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction data"})
		return
	}

	if err := tx.Where("transaction_id = ?", existingTransaction.ID).Delete(&model.Delivery{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old delivery data"})
		return
	}

	for _, delivery := range request.Deliveries {
		delivery.ID = 0
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

func (c *TransactionController) DeleteTransaction(ctx *gin.Context) {
	id := ctx.Param("id")

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

	if err := tx.Where("transaction_id = ?", id).Delete(&model.Delivery{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete delivery data"})
		return
	}

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

func (c *TransactionController) PatchTransaction(ctx *gin.Context) {
	id := ctx.Param("id")

	var request struct {
		DownPayment       *int       `json:"down_payment"`
		DownPaymentStatus *string    `json:"down_payment_status"`
		DownPaymentTime   *time.Time `json:"down_payment_time"`
		FullPayment       *float64   `json:"full_payment"`
		FullPaymentStatus *string    `json:"full_payment_status"`
		FullPaymentTime   *time.Time `json:"full_payment_time"`
		TransactionStatus *string    `json:"transaction_status"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var transaction model.Transaction
	if err := c.DB.First(&transaction, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data transaksi"})
		return
	}

	updates := map[string]interface{}{}

	if request.DownPayment != nil {
		updates["down_payment"] = *request.DownPayment
	}
	if request.DownPaymentTime != nil {
		updates["down_payment_time"] = *request.DownPaymentTime
	}
	if request.DownPaymentStatus != nil {
		updates["down_payment_status"] = *request.DownPaymentStatus
	}

	if request.FullPayment != nil {
		updates["full_payment"] = *request.FullPayment
	}
	if request.FullPaymentStatus != nil {
		updates["full_payment_status"] = *request.FullPaymentStatus
	}
	if request.FullPaymentTime != nil {
		updates["full_payment_time"] = *request.FullPaymentTime
	}

	if request.TransactionStatus != nil {
		updates["transaction_status"] = *request.TransactionStatus
	}

	if len(updates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tidak ada data yang diperbarui"})
		return
	}

	if err := c.DB.Model(&transaction).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate transaksi"})
		return
	}

	if err := c.DB.Preload("Customer").Preload("Delivery").First(&transaction, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data terbaru"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transaksi berhasil diperbarui",
		"data":    transaction,
	})
}

func (c *TransactionController) UpdateTransactionStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	var body struct {
		TransactionStatus string `json:"transaction_status" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	tx := c.DB.Begin()

	var transaction model.Transaction
	if err := tx.First(&transaction, id).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}

	transaction.TransactionStatus = body.TransactionStatus
	if err := tx.Save(&transaction).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate transaksi"})
		return
	}

	if body.TransactionStatus == "berjalan" {
		if err := tx.Model(&model.Delivery{}).
			Where("transaction_id = ? AND delivery_status = ?", transaction.ID, "disetujui").
			Update("delivery_status", "menunggu pengemudi").Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status pengiriman"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Status transaksi (dan pengiriman jika perlu) berhasil diperbarui",
		"data":    transaction,
	})
}
