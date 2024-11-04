package transactioncontroller

import (
	"EcommerceSederhana/config"
	"EcommerceSederhana/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllTransaction(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
		return
	}

	var transaction []models.Transaction

	if err := config.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Preload("TransactionItem.Product").Where("user_id =?", userId).Find(&transaction).Error; err != nil {
		c.JSON(404, gin.H{"message": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transaction})
}

func GetTransactionById(c *gin.Context) {
	transactionId := c.Param("transactionId")
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
		return
	}

	var transaction models.Transaction

	if err := config.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Preload("TransactionItem.Product").Where("id = ?", transactionId).Where("user_id =?", userId).First(&transaction).Error; err != nil {
		c.JSON(404, gin.H{"message": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

func Checkout(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
		return
	}

	var cart []models.Cart
	var totalPrice float64

	// Ambil cart
	if err := config.DB.Where("user_id = ?", userId).Find(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if len(cart) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "cart is empty"})
		return
	}

	// Cek produk pada setiap cart apakah stock ada
	for _, item := range cart {
		var product models.Product

		if err := config.DB.Where("id = ?", item.ProductId).First(&product).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
			return
		}

		if product.Stock < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Product '%s' out of stock", product.Name),
			})
			return
		}

		totalPrice += item.TotalPrice
	}

	// Buat Transaksi baru
	var transctionreq models.TransactionReq
	if err := c.ShouldBind(&transctionreq); err != nil {
		c.JSON(422, gin.H{"message": err.Error()})
		return
	}

	transaction := models.Transaction{
		UserId:          userId.(int),
		TotalPrice:      totalPrice,
		ShippingAddress: transctionreq.ShippingAddress,
		PaymentStatus:   models.Status(models.Pending),
	}

	tx := config.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Buat Transaction Item dan kurangi stock
	for _, item := range cart {
		var product models.Product

		if err := tx.Where("id = ?", item.ProductId).First(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
			return
		}

		transactionItem := models.TransactionItem{
			ProductId:     item.ProductId,
			TransactionId: transaction.Id,
			Quantity:      item.Quantity,
			TotalPrice:    item.TotalPrice,
		}

		if err := tx.Create(&transactionItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		product.Stock -= item.Quantity

		if product.Stock < 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Insufficient stock for product",
			})
			return
		}

		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	// Kosongkan cart
	if result := tx.Where("user_id = ?", userId).Delete(&models.Cart{}); result.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Can't delete cart"})
		return
	}

	// Commit jika semua berhasil tanpa ada kesalahan
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to commit transaction",
			"error":   err.Error(),
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message":     "Checkout successful",
		"transaction": transaction,
	})
}

func CancelTransaction(c *gin.Context) {
	transactionId := c.Param("transactionId")
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
		return
	}

	tx := config.DB.Begin()

	if result := tx.Model(&models.TransactionItem{}).Where("transaction_id =?", transactionId).Delete(&models.TransactionItem{}); result.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Can't delete transaction item"})
		return
	}

	if result := tx.Model(&models.Transaction{}).Where("user_id = ?", userId).Where("id =?", transactionId).Delete(&models.Transaction{}); result.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Can't delete transaction"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to commit transaction: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cancel Transaction success",
	})
}
