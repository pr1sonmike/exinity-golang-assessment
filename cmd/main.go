package main

import (
	"exinity-golang-assessment/config"
	"exinity-golang-assessment/internal/controllers"
	"exinity-golang-assessment/internal/gateways"
	"exinity-golang-assessment/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	r := gin.Default()
	db := config.InitDB()

	gatewayA := gateways.NewGatewayA()
	gatewayB := gateways.NewGatewayB()

	transactionService := services.NewTransactionService(db, gatewayA, gatewayB)

	// Routes
	r.POST("/deposit", func(c *gin.Context) {
		controllers.Deposit(c, transactionService)
	})
	r.POST("/withdraw", func(c *gin.Context) {
		controllers.Withdraw(c, transactionService)
	})
	r.POST("/callback", func(c *gin.Context) {
		controllers.HandleCallback(c, transactionService)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
