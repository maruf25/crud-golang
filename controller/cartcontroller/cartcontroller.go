package cartcontroller

import (
	"EcommerceSederhana/config"
	"EcommerceSederhana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllCart(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	var cart []models.Cart

	if err := config.DB.Preload("Product").Where("user_id = ?", userId).Find(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve cart",
			"data":    err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Get All Cart",
		"cart":    cart,
	})
}

func AddToCart(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
		return
	}

	var createCart models.CreateCart
	var product models.Product

	if err := c.ShouldBindJSON(&createCart); err != nil {
		c.AbortWithStatusJSON(422, gin.H{
			"message": "Validation error",
			"data":    err.Error(),
		})
		return
	}

	if err := config.DB.Where("id = ?", createCart.ProductId).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "product not found",
		})
		return
	}

	cart := &models.Cart{
		ProductId:  product.Id,
		Quantity:   createCart.Quantity,
		UserId:     userId.(int),
		TotalPrice: product.Price * float64(createCart.Quantity),
	}

	// Cek apakah product sudah pernah dimasukkan kedalam cart
	err := config.DB.Where("user_id = ?", userId).Where("product_id = ?", product.Id).First(&cart).Error
	if err == nil {
		cart.Quantity += createCart.Quantity
		cart.TotalPrice = product.Price * float64(cart.Quantity)

		err := config.DB.Save(&cart).Error
		// Ketika error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed add item to cart",
				"data":    err.Error(),
			})
			return
		}
		// // Ketika berhasil
		// c.JSON(200, gin.H{
		// 	"message": "Success add to cart",
		// 	"cart":    existingItemInCart,
		// })
		// return
	} else if err := config.DB.Create(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Tampilkan Data yang baru saja dimasukkan
	if err := config.DB.Preload("Product").First(&cart, cart.Id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve cart",
			"data":    err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Success add to cart",
		"cart":    cart,
	})
}

func RemoveFromCart(c *gin.Context) {
	cartId := c.Param("cartId")
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
		return
	}

	var cart models.Cart

	if err := config.DB.Where("id = ?", cartId).Where("user_id = ?", userId).First(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to remove item from cart",
			"data":    err.Error(),
		})
		return
	}

	if err := config.DB.Delete(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to remove item from cart",
			"data":    err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Success remove item from cart",
	})
}
