package controllers

import (
	"ady-trans-jaya-golang/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func generateDeliveryCode(db *gorm.DB) (string, error) {
	today := time.Now().Format("20060102")
	prefix := "DLV-" + today

	var count int64
	if err := db.Model(&model.Delivery{}).
		Where("to_char(created_at, 'YYYYMMDD') = ?", today).
		Count(&count).Error; err != nil {
		return "", err
	}

	serial := fmt.Sprintf("%04d", count+1)
	return prefix + "-" + serial, nil
}

func DeliveryControllers(r *gin.Engine, db *gorm.DB) {
	r.GET("/api/deliveries", func(ctx *gin.Context) {
		var delivery []model.Delivery
		if err := db.Preload("Transaction").Preload("Transaction.Customer").Preload("Driver").Preload("Vehicle").Find(&delivery).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve delivery data"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": delivery})
	})

	r.POST("/api/deliveries", func(ctx *gin.Context) {
		var delivery model.Delivery

		if err := ctx.ShouldBindJSON(&delivery); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		code, err := generateDeliveryCode(db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate delivery code"})
			return
		}
		delivery.DeliveryCode = code

		if delivery.DeliveryStatus == "" {
			delivery.DeliveryStatus = "menunggu persetujuan"
		}

		// Simpan ke database
		if err := db.Create(&delivery).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save delivery data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Delivery data successfully saved", "data": delivery})
	})

	r.GET("/api/delivery/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var delivery model.Delivery

		if err := db.Preload("Transaction").Preload("Transaction.Customer").Preload("Driver").Preload("Vehicle").First(&delivery, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Delivery data not found"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": delivery})
	})

	r.GET("/api/deliveries/driver/:driverId", func(ctx *gin.Context) {
		id := ctx.Param("driverId")

		var deliveries []model.Delivery
		if err := db.Preload("Transaction").Preload("Transaction.Customer").
			Preload("Driver").Preload("Vehicle").
			Where("driver_id = ?", id).
			Find(&deliveries).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deliveries"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": deliveries})
	})

	// Search deliveries with query parameters
	r.GET("/api/delivery/search", func(ctx *gin.Context) {
		var deliveries []model.Delivery
		query := db.Preload("Transaction").Preload("Transaction.Customer").Preload("Driver").Preload("Vehicle")

		// Process query parameters
		if status := ctx.Query("status"); status != "" {
			query = query.Where("delivery_status = ?", status)
		}

		if driverID := ctx.Query("driver_id"); driverID != "" {
			query = query.Where("driver_id = ?", driverID)
		}

		if vehicleID := ctx.Query("vehicle_id"); vehicleID != "" {
			query = query.Where("vehicle_id = ?", vehicleID)
		}

		if transactionID := ctx.Query("transaction_id"); transactionID != "" {
			query = query.Where("transaction_id = ?", transactionID)
		}

		// Date range searching if needed
		if startDate := ctx.Query("start_date"); startDate != "" {
			query = query.Where("created_at >= ?", startDate)
		}

		if endDate := ctx.Query("end_date"); endDate != "" {
			query = query.Where("created_at <= ?", endDate)
		}

		// Execute the query
		if err := query.Find(&deliveries).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search delivery data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": deliveries})
	})

	// Update/Edit delivery
	r.PUT("/api/delivery/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var delivery model.Delivery

		// Check if delivery exists
		if err := db.First(&delivery, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Delivery data not found"})
			return
		}

		// Bind the updated data
		var updatedDelivery model.Delivery
		if err := ctx.ShouldBindJSON(&updatedDelivery); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the delivery
		if err := db.Model(&delivery).Updates(updatedDelivery).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update delivery data"})
			return
		}

		// Fetch the updated delivery with all relations
		if err := db.Preload("Transaction").Preload("Transaction.Customer").Preload("Driver").Preload("Vehicle").First(&delivery, id).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated delivery data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Delivery data successfully updated", "data": delivery})
	})

	// Delete delivery
	r.DELETE("/api/delivery/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var delivery model.Delivery

		// Check if delivery exists
		if err := db.First(&delivery, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Delivery data not found"})
			return
		}

		// Delete the delivery
		if err := db.Delete(&delivery).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete delivery data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Delivery data successfully deleted"})
	})

	// Update delivery status specifically
	r.PATCH("/api/delivery/:id/status", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var delivery model.Delivery

		// Check if delivery exists
		if err := db.First(&delivery, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Delivery data not found"})
			return
		}

		// Bind only the status update
		var statusUpdate struct {
			DeliveryStatus string `json:"delivery_status"`
		}

		if err := ctx.ShouldBindJSON(&statusUpdate); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update only the delivery status
		if err := db.Model(&delivery).Update("delivery_status", statusUpdate.DeliveryStatus).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update delivery status"})
			return
		}

		// Fetch the updated delivery
		if err := db.Preload("Transaction").Preload("Transaction.Customer").Preload("Driver").Preload("Vehicle").First(&delivery, id).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated delivery data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Delivery status successfully updated", "data": delivery})
	})
}
