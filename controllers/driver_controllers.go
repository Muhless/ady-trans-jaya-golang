package controllers

import (
	"ady-trans-jaya-golang/model"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func DriversControllers(r *gin.Engine, db *gorm.DB) {
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	// GET semua driver
	r.GET("/api/driver", func(ctx *gin.Context) {
		var driver []model.Driver

		// Pencarian berdasarkan nama atau nomor telepon
		search := ctx.Query("search")
		if search != "" {
			if err := db.Where("name LIKE ? OR phone LIKE ?", "%"+search+"%", "%"+search+"%").Find(&driver).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search drivers data"})
				return
			}
		} else {
			if err := db.Find(&driver).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve drivers data"})
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"data": driver})
	})

	// GET driver by ID
	r.GET("api/driver/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var driver model.Driver
		if err := db.First(&driver, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": driver})
	})

	r.POST("/api/driver", func(ctx *gin.Context) {
		var driver model.Driver

		if err := ctx.ShouldBind(&driver); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		file, err := ctx.FormFile("photo")
		if err == nil {
			filename := uuid.New().String() + filepath.Ext(file.Filename)
			filePath := filepath.Join(uploadDir, filename)

			if err := ctx.SaveUploadedFile(file, filePath); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
				return
			}

			driver.Photo = "/uploads/" + filename
		}

		if driver.Status == "" {
			driver.Status = "tersedia"
		}

		var existingDriver model.Driver
		if err := db.Where("phone = ?", driver.Phone).First(&existingDriver).Error; err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nomor telepon sudah terdaftar"})
			return
		}

		if err := db.Create(&driver).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed save drivers data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Driver successfully saved", "data": driver})
	})

	r.PUT("/api/driver/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		var existingDriver model.Driver
		if err := db.First(&existingDriver, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
			return
		}

		var updatedDriver model.Driver
		if err := ctx.ShouldBind(&updatedDriver); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		file, err := ctx.FormFile("photo")
		if err == nil {
			if existingDriver.Photo != "" && !strings.HasPrefix(existingDriver.Photo, "http") {
				oldFilePath := "." + existingDriver.Photo
				if _, err := os.Stat(oldFilePath); err == nil {
					os.Remove(oldFilePath)
				}
			}

			filename := uuid.New().String() + filepath.Ext(file.Filename)
			filePath := filepath.Join(uploadDir, filename)

			if err := ctx.SaveUploadedFile(file, filePath); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
				return
			}

			updatedDriver.Photo = "/uploads/" + filename
		} else {
			updatedDriver.Photo = existingDriver.Photo
		}

		if updatedDriver.Phone != existingDriver.Phone {
			var phoneExist model.Driver
			if err := db.Where("phone = ? AND id != ?", updatedDriver.Phone, id).First(&phoneExist).Error; err == nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nomor telepon sudah terdaftar"})
				return
			}
		}

		updatedDriver.ID = existingDriver.ID
		updatedDriver.CreatedAt = existingDriver.CreatedAt
		updatedDriver.UpdatedAt = time.Now()

		if err := db.Save(&updatedDriver).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update driver"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Driver successfully updated", "data": updatedDriver})
	})

	r.DELETE("/api/driver/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		var driver model.Driver
		if err := db.First(&driver, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
			return
		}

		if driver.Photo != "" && !strings.HasPrefix(driver.Photo, "http") {
			filePath := "." + driver.Photo
			if _, err := os.Stat(filePath); err == nil {
				os.Remove(filePath)
			}
		}

		if err := db.Delete(&driver).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete driver"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Driver successfully deleted"})
	})

	r.GET("/api/driver/search", func(ctx *gin.Context) {
		search := ctx.Query("q")
		if search == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
			return
		}

		var drivers []model.Driver
		if err := db.Where("name LIKE ? OR phone LIKE ? OR address LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%").Find(&drivers).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search drivers"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": drivers})
	})

	r.Static("/uploads", "./uploads")
}
