package productcontroller

import (
	"EcommerceSederhana/config"
	"EcommerceSederhana/models"
	"EcommerceSederhana/utils"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(c *gin.Context) {
	var products []models.Product

	if err := config.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"products": products,
	})
}

func GetProductById(c *gin.Context) {
	id := c.Param("id")

	var product models.Product

	if err := config.DB.Where("id = ?", id).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "product not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"product": product,
	})
}

func CreateProduct(c *gin.Context) {
	var product models.Product

	image, _ := c.FormFile("File")

	// Cek image harus ada
	if image == nil {
		c.AbortWithStatusJSON(422, gin.H{
			"message": "Validation error",
			"data":    "You should upload image",
		})
		return
	}

	fileType := image.Header.Get("Content-Type")

	// Cek mimetype harus jpg,jpeg, atau png
	if !utils.IsSupportedFileType(fileType) {
		c.AbortWithStatusJSON(422, gin.H{
			"message": "Validation error",
			"data":    "You should upload image with extension jpg, jpeg, or png",
		})
		return
	}

	fileName := time.Now().Format("2006-01-02") + image.Filename
	currentDir, _ := os.Getwd()

	imagePath := filepath.Join(currentDir, "images/", fileName)

	err := c.SaveUploadedFile(image, imagePath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Validation error",
			"data":    err.Error(),
		})
		return
	}

	product.Image = "images/" + fileName

	if err := c.ShouldBind(&product); err != nil {
		c.AbortWithStatusJSON(422, gin.H{
			"message": "Validation error",
			"data":    err.Error(),
		})
		return
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "Product succes create",
		"product": product,
	})
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product

	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	image, _ := c.FormFile("File")

	// Cek image jika melakukan update maka
	if image != nil {
		fileType := image.Header.Get("Content-Type")

		// Cek mimetype harus jpg,jpeg, atau png
		if !utils.IsSupportedFileType(fileType) {
			c.AbortWithStatusJSON(422, gin.H{
				"message": "Validation error",
				"data":    "You should upload image with extension jpg, jpeg, or png",
			})
			return
		}

		fileName := time.Now().Format("2006-01-02") + image.Filename
		currentDir, _ := os.Getwd()

		imagePath := filepath.Join(currentDir, "images/", fileName)

		err := c.SaveUploadedFile(image, imagePath)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Validation error",
				"data":    err.Error(),
			})
			return
		}

		imagePath = filepath.Join(currentDir, product.Image)
		errDelete := os.Remove(imagePath)
		if errDelete != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Validation error",
				"data":    errDelete.Error(),
			})
			return
		}

		product.Image = "images/" + fileName
	}

	if err := c.ShouldBind(&product); err != nil {
		c.AbortWithStatusJSON(422, gin.H{
			"message": "Validation error",
			"data":    err.Error(),
		})
		return
	}

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Product succes update",
		"product": product,
	})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Product succes Delete",
	})
}
