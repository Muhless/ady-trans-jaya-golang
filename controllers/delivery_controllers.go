package controllers

import (
	"ady-trans-jaya-golang/db"
	"ady-trans-jaya-golang/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeliveryController struct {
	DB *gorm.DB
}

func DeliveryControllers(r *gin.Engine, db *gorm.DB) {
	r.GET("/api/deliveries", func(ctx *gin.Context) {
		var delivery []model.Delivery
		if err := db.Preload("Transaction").
			Preload("Transaction.Customer").
			Preload("Driver").
			Preload("Vehicle").
			Preload("Items").
			Preload("DeliveryProgress").
			Find(&delivery).Error; err != nil {
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

		if err := db.Create(&delivery).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save delivery data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Delivery data successfully saved", "data": delivery})
	})

	r.GET("/api/delivery/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var delivery model.Delivery

		if err := db.Preload("Transaction").Preload("Transaction.Customer").Preload("Driver").Preload("Vehicle").Preload("Items").Preload("DeliveryProgress").First(&delivery, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Delivery data not found"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": delivery})
	})

	r.GET("/api/deliveries/driver/:driverId", func(ctx *gin.Context) {
		id := ctx.Param("driverId")

		var deliveries []model.Delivery
		if err := db.Preload("Transaction").Preload("Transaction.Customer").
			Preload("Driver").Preload("Vehicle").Preload("Items").Preload("DeliveryProgress").
			Where("driver_id = ?", id).
			Find(&deliveries).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deliveries"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": deliveries})
	})

	r.GET("/api/delivery/search", func(ctx *gin.Context) {
		var deliveries []model.Delivery
		query := db.Preload("Transaction").Preload("Transaction.Customer").Preload("Driver").Preload("Vehicle")

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

		if startDate := ctx.Query("start_date"); startDate != "" {
			query = query.Where("created_at >= ?", startDate)
		}

		if endDate := ctx.Query("end_date"); endDate != "" {
			query = query.Where("created_at <= ?", endDate)
		}

		if err := query.Find(&deliveries).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search delivery data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": deliveries})
	})

	r.PUT("/api/delivery/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var delivery model.Delivery

		if err := db.First(&delivery, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Delivery data not found"})
			return
		}

		var updatedDelivery model.Delivery
		if err := ctx.ShouldBindJSON(&updatedDelivery); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Model(&delivery).Updates(updatedDelivery).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update delivery data"})
			return
		}

		if err := db.Preload("Transaction").Preload("Transaction.Customer").Preload("Driver").Preload("Vehicle").First(&delivery, id).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated delivery data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Delivery data successfully updated", "data": delivery})
	})

	r.PATCH("/api/delivery/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var delivery model.Delivery

		if err := db.First(&delivery, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Delivery data not found"})
			return
		}

		var updateFields map[string]interface{}
		if err := ctx.ShouldBindJSON(&updateFields); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Model(&delivery).Updates(updateFields).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update delivery data"})
			return
		}

		if err := db.Preload("Transaction").Preload("Transaction.Customer").Preload("Driver").Preload("Vehicle").First(&delivery, id).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated delivery data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Delivery data successfully updated",
			"data":    delivery,
		})
	})

	r.DELETE("/api/delivery/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var delivery model.Delivery

		if err := db.First(&delivery, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Delivery data not found"})
			return
		}

		if err := db.Delete(&delivery).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete delivery data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Delivery data successfully deleted"})
	})

}

func generateDeliveryCode(db *gorm.DB) (string, error) {
	today := time.Now().Format("20060102")
	prefix := "ATJ-" + today

	var count int64
	if err := db.Model(&model.Delivery{}).
		Where("to_char(created_at, 'YYYYMMDD') = ?", today).
		Count(&count).Error; err != nil {
		return "", err
	}

	serial := fmt.Sprintf("%04d", count+1)
	return prefix + "-" + serial, nil
}

func (c *DeliveryController) UpdateDeliveryStatus(ctx *gin.Context) {
	var body struct {
		DeliveryStatus string `json:"delivery_status" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	validStatuses := []string{"menunggu persetujuan", "disetujui", "ditolak", "dalam pengiriman"}
	isValid := false
	for _, status := range validStatuses {
		if body.DeliveryStatus == status {
			isValid = true
			break
		}
	}
	if !isValid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Status pengiriman tidak valid"})
		return
	}

	id := ctx.Param("id")

	tx := c.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var delivery model.Delivery
	if err := tx.First(&delivery, id).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Pengiriman tidak ditemukan"})
		return
	}

	delivery.DeliveryStatus = model.DeliveryStatus(body.DeliveryStatus)
	if body.DeliveryStatus == "disetujui" {
		now := time.Now()
		delivery.ApprovedAt = &now
	}

	if err := tx.Save(&delivery).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status pengiriman"})
		return
	}

	if body.DeliveryStatus == "disetujui" || body.DeliveryStatus == "ditolak" {
		var deliveries []model.Delivery
		if err := tx.Where("transaction_id = ?", delivery.TransactionID).Find(&deliveries).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar pengiriman"})
			return
		}

		approvedCount := 0
		rejectedCount := 0

		for _, d := range deliveries {
			if string(d.DeliveryStatus) == "disetujui" {
				approvedCount++
			} else if string(d.DeliveryStatus) == "ditolak" {
				rejectedCount++
			}
		}

		totalDeliveries := len(deliveries)

		if rejectedCount == totalDeliveries {
			if err := tx.Model(&model.Transaction{}).
				Where("id = ?", delivery.TransactionID).
				Update("transaction_status", "dibatalkan").Error; err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status transaksi"})
				return
			}
		} else if approvedCount >= 1 {
			if err := tx.Model(&model.Transaction{}).
				Where("id = ?", delivery.TransactionID).
				Update("transaction_status", "diproses").Error; err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status transaksi"})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Status pengiriman berhasil diperbarui",
		"data":    delivery,
	})
}

func (c *DeliveryController) GetHistoryDeliveries(ctx *gin.Context) {
	driverID := ctx.Param("id")

	var deliveries []model.Delivery
	if err := c.DB.Preload("Transaction").
		Preload("Transaction.Customer").
		Preload("Driver").
		Preload("Vehicle").
		Preload("Items").
		Preload("DeliveryProgress").
		Where("driver_id = ? AND (delivery_status = ? OR delivery_status = ?)",
			driverID,
			"selesai",
			"dibatalkan",
		).
		Find(&deliveries).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengiriman"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": deliveries})
}

func GetActiveDeliveriesByDriver(ctx *gin.Context) {
	driverID := ctx.Param("id")

	var deliveries []model.Delivery
	if err := db.
		DB.Preload("Transaction").
		Preload("Transaction.Customer").
		Preload("Driver").
		Preload("Vehicle").
		Preload("Items").
		Preload("DeliveryProgress").
		Where("driver_id = ? AND (delivery_status = ? OR delivery_status = ?)",
			driverID,
			model.DeliveryStatusOnDelivery,
			model.DeliveryStatusWaitingDriver,
		).
		Find(&deliveries).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil pengiriman aktif"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": deliveries})
}

func (c *DeliveryController) UpdateDeliveryDriverStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	var req model.Delivery

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	var delivery model.Delivery
	if err := c.DB.First(&delivery, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	delivery.DeliveryStatus = req.DeliveryStatus

	if err := c.DB.Save(&delivery).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	if req.DeliveryStatus == "selesai" {
		if err := c.DB.Model(&model.Driver{}).
			Where("id = ?", delivery.DriverID).
			Update("status", "tersedia").Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status driver"})
			return
		}

		if err := c.DB.Model(&model.Vehicle{}).
			Where("id = ?", delivery.VehicleID).
			Update("status", "tersedia").Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status kendaraan"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Status pengiriman berhasil diperbarui"})
}
