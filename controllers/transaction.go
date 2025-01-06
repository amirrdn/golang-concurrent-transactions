package controllers

import (
	"e-wallet/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var walletService *services.EWalletSystem

func InitWalletService(ws *services.EWalletSystem) {
	walletService = ws
}

func Credit(c *gin.Context) {
	var request struct {
		UserID int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	balance, transactionID, err := walletService.Credit(request.UserID, request.Amount)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"transaction_id": transactionID,
		"new_balance":    balance,
	})
}

func Debit(c *gin.Context) {
	var request struct {
		UserID int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	balance, transactionID, err := walletService.Debit(request.UserID, request.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"transaction_id": transactionID,
		"new_balance":    balance,
	})
}
