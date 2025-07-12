package controllers

import (
	"ady-trans-jaya-golang/db"
	"ady-trans-jaya-golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CustomersController struct {
	DB *gorm.DB
}

func CustomersControllers(r *gin.Engine, db *gorm.DB) {
	r.GET("/api/customers", func(ctx *gin.Context) {
		var customer []model.Customer
		if err := db.Order("id ASC").Find(&customer).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers data"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": customer})
	})

	r.GET("/api/customer/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var customer model.Customer
		if err := db.First(&customer, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": customer})
	})

	r.POST("/api/customer", func(ctx *gin.Context) {
		var customer model.Customer

		if err := ctx.ShouldBindJSON(&customer); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&customer).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save customer data"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Customer data successfully saved", "data": customer})
	})

}

func (c *CustomersController) UpdateCustomer(ctx *gin.Context) {

	id := ctx.Param("id")

	var customer model.Customer
	if err := db.DB.First(&customer, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	var input struct {
		Name    string `json:"name"`
		Company string `json:"company"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer.Name = input.Name
	customer.Company = input.Company
	customer.Email = input.Email
	customer.Phone = input.Phone
	customer.Address = input.Address

	if err := db.DB.Save(&customer).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui pelanggan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Data pelanggan berhasil diperbarui", "data": customer})
}
