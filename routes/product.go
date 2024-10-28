package routes

import (
	"EcommerceSederhana/controller/productcontroller"
	"EcommerceSederhana/middlewares"
	"EcommerceSederhana/models"

	"github.com/gin-gonic/gin"
)

func ProductRoute(r *gin.Engine) {
	productGroup := r.Group("/products")
	productGroup.GET("/", productcontroller.GetAllProducts)
	productGroup.GET("/:id", productcontroller.GetProductById)
	productGroup.Use(middlewares.AuthMiddleware([]models.Role{models.Admin}))
	{
		productGroup.POST("/", productcontroller.CreateProduct)
		productGroup.PUT("/:id", productcontroller.UpdateProduct)
		productGroup.DELETE("/:id", productcontroller.DeleteProduct)
	}
}
