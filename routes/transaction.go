package routes

import (
	"EcommerceSederhana/controller/transactioncontroller"
	"EcommerceSederhana/middlewares"
	"EcommerceSederhana/models"

	"github.com/gin-gonic/gin"
)

func TransactionRoute(r *gin.Engine) {
	transactionGroup := r.Group("/transactions")

	transactionGroup.Use(middlewares.AuthMiddleware([]models.Role{models.Admin, models.Member}))
	{
		transactionGroup.GET("/", transactioncontroller.GetAllTransaction)
		transactionGroup.GET("/:transactionId", transactioncontroller.GetTransactionById)
		transactionGroup.POST("/", transactioncontroller.Checkout)
		transactionGroup.DELETE("/:transactionId", transactioncontroller.CancelTransaction)
	}
}
