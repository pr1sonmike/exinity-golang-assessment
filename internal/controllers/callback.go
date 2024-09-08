package controllers

import (
	"exinity-golang-assessment/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleCallback(c *gin.Context, service *services.TransactionService) {
	var req services.CallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid callback data"})
		return
	}

	if err := service.HandleCallback(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Transaction updated"})
}
