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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func DriversControllers(r *gin.Engine, db *gorm.DB) {
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	r.GET("/api/drivers", func(ctx *gin.Context) {
		var driver []model.Driver
		db.Order("created_at DESC").Find(&driver)

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
		// Parse multipart form
		if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Form parsing failed: " + err.Error()})
			return
		}

		// Ambil nilai dari form
		name := ctx.PostForm("name")
		phone := ctx.PostForm("phone")
		address := ctx.PostForm("address")
		status := ctx.PostForm("status")
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")

		if name == "" || phone == "" || username == "" || password == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nama, Telepon, Username, dan Password wajib diisi"})
			return
		}

		// Cek duplikasi phone
		var exist model.Driver
		if err := db.Where("phone = ?", phone).Or("username = ?", username).First(&exist).Error; err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username atau nomor telepon sudah digunakan"})
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password"})
			return
		}

		// Handle upload photo
		var photoPath string
		file, err := ctx.FormFile("photo")
		if err == nil {
			filename := uuid.New().String() + filepath.Ext(file.Filename)
			filePath := filepath.Join(uploadDir, filename)
			if err := ctx.SaveUploadedFile(file, filePath); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file foto"})
				return
			}
			photoPath = "/uploads/" + filename
		}

		// Default status jika kosong
		if status == "" {
			status = "tersedia"
		}

		// Buat driver
		driver := model.Driver{
			Name:     name,
			Phone:    phone,
			Address:  address,
			Status:   model.DriverStatus(status),
			Username: username,
			Password: string(hashedPassword),
			Photo:    photoPath,
		}

		if err := db.Create(&driver).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data driver"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Driver berhasil disimpan",
			"data":    driver,
		})
	})

	r.PUT("/api/driver/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		var existingDriver model.Driver
		if err := db.First(&existingDriver, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
			return
		}

		if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
			return
		}

		var updatedDriver model.Driver
		updatedDriver.Name = ctx.PostForm("name")
		updatedDriver.Phone = ctx.PostForm("phone")
		updatedDriver.Address = ctx.PostForm("address")
		updatedDriver.Status = model.DriverStatus(ctx.PostForm("status"))
		updatedDriver.Username = ctx.PostForm("username")
		password := ctx.PostForm("password")
		if password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
				return
			}
			existingDriver.Password = string(hashedPassword)
		}

		// Handle uploaded file
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

		// Cek apakah phone sudah digunakan driver lain
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

	r.PATCH("/api/driver/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		var payload map[string]interface{}
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		// Validasi status jika ada
		if statusRaw, ok := payload["status"]; ok {
			status, ok := statusRaw.(string)
			if !ok || (status != "tersedia" && status != "tidak tersedia") {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Status harus 'tersedia' atau 'tidak tersedia'"})
				return
			}
		}

		// Hash password jika akan diubah
		if passwordRaw, ok := payload["password"]; ok {
			passwordStr, ok := passwordRaw.(string)
			if !ok || passwordStr == "" {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password tidak boleh kosong"})
				return
			}
			hashed, err := bcrypt.GenerateFromPassword([]byte(passwordStr), bcrypt.DefaultCost)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hash password"})
				return
			}
			payload["password"] = string(hashed)
		}

		payload["updated_at"] = time.Now()

		if err := db.Model(&model.Driver{}).Where("id = ?", id).Updates(payload).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate data driver"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Driver berhasil diupdate"})
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
