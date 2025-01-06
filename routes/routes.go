package routes

import (
	"e-wallet/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/api/transactions/credit", controllers.Credit)
	r.POST("/api/transactions/debit", controllers.Debit)
}
