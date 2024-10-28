package routes

import (
	"EcommerceSederhana/controller/authcontroller"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/")
	{
		authGroup.POST("/register", authcontroller.Register)
		authGroup.POST("/login", authcontroller.Login)
	}
}
