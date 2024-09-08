package controllers

import (
	"exinity-golang-assessment/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Deposit(c *gin.Context, service *services.TransactionService) {
	var req services.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	res, err := service.Deposit(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func Withdraw(c *gin.Context, service *services.TransactionService) {
	var req services.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	res, err := service.Withdraw(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
