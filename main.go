package main

import (
	"EcommerceSederhana/config"
	"EcommerceSederhana/middlewares"
	"EcommerceSederhana/models"
	"EcommerceSederhana/routes"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/welcome", middlewares.AuthMiddleware([]models.Role{models.Member, models.Admin}), func(c *gin.Context) {
		email, _ := c.Get("email")
		c.JSON(http.StatusOK, gin.H{"message": "selamat datang " + email.(string)})
	})

	// Menampilkan Gambar
	r.GET("/images/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		imagePath := filepath.Join("images", filename)

		c.File(imagePath)
	})

	routes.AuthRoutes(r)
	routes.ProductRoute(r)

	config.ConnectDB()

	r.Run(":8000")
}
