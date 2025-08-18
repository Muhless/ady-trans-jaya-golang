package controllers

import (
	"ady-trans-jaya-golang/model"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func VehicleControllers(r *gin.Engine, db *gorm.DB) {
	// fungsi get data
	r.GET("/api/vehicles", func(ctx *gin.Context) {
		var vehicle []model.Vehicle
		if err := db.Order("created_at ASC").Find(&vehicle).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kendaraan"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": vehicle})
	})

	r.GET("/api/vehicle/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var vehicle model.Vehicle
		if err := db.First(&vehicle, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Kendaraan tidak ditemukan"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": vehicle})
	})

	// fungsi create data
	r.POST("/api/vehicle", func(ctx *gin.Context) {
		var vehicle model.Vehicle

		body, _ := io.ReadAll(ctx.Request.Body)
		fmt.Println("📥 Raw Body:", string(body))
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		if err := ctx.ShouldBindJSON(&vehicle); err != nil {
			fmt.Println("❌ Error binding:", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if vehicle.Status == "" {
			vehicle.Status = "tersedia"
		}

		if err := db.Create(&vehicle).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data kendaraan"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Data kendaraan berhasil disimpan", "data": vehicle})
	})

	r.PUT("/api/vehicle/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var vehicle model.Vehicle

		if err := db.First(&vehicle, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			return
		}

		if err := ctx.ShouldBindJSON(&vehicle); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&vehicle).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah data kendaraan"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Data kendaraan berhasil diperbarui", "data": vehicle})
	})

	r.PATCH("/api/vehicle/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		var payload map[string]interface{}
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		if statusRaw, ok := payload["status"]; ok {
			status, ok := statusRaw.(string)
			if !ok || (status != "tersedia" && status != "tidak tersedia") {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Status harus 'tersedia' atau 'tidak tersedia'"})
				return
			}
		}

		if err := db.Model(&model.Vehicle{}).Where("id = ?", id).Updates(payload).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Data Kendaraan berhasil diperbarui", "data": payload})
	})

	r.DELETE("/api/vehicle/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var vehicle model.Vehicle

		if err := db.First(&vehicle, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			return
		}

		if err := db.Delete(&vehicle).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data kendaraan"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Data kendaraan berhasil dihapus"})
	})

	r.GET("/api/vehicle/search", func(ctx *gin.Context) {
		searchQuery := ctx.DefaultQuery("query", "")
		if searchQuery == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
			return
		}

		var vehicles []model.Vehicle
		if err := db.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(searchQuery)+"%").
			Or("LOWER(license_plate) LIKE ?", "%"+strings.ToLower(searchQuery)+"%").
			Find(&vehicles).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search vehicles"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": vehicles})
	})
}
