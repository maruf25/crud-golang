package routes

import (
	"EcommerceSederhana/controller/cartcontroller"
	"EcommerceSederhana/middlewares"
	"EcommerceSederhana/models"

	"github.com/gin-gonic/gin"
)

func CartRoute(r *gin.Engine) {
	cartGroupe := r.Group("/carts")

	cartGroupe.Use(middlewares.AuthMiddleware([]models.Role{models.Admin, models.Member}))
	{
		cartGroupe.GET("/", cartcontroller.GetAllCart)
		cartGroupe.POST("/", cartcontroller.AddToCart)
		cartGroupe.DELETE("/:cartId", cartcontroller.RemoveFromCart)
	}
}
