package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"exinity-golang-assessment/config"
	"exinity-golang-assessment/internal/gateways"
	"exinity-golang-assessment/internal/handlers"
	"exinity-golang-assessment/internal/mocks"
	"exinity-golang-assessment/internal/repository"
	"exinity-golang-assessment/internal/service"
	"exinity-golang-assessment/internal/utils"
)

func main() {
	utils.InitLogger()

	cfg, err := config.LoadConfig("config")
	if err != nil {
		utils.Logger.Fatalf("Failed to load config: %v", err)
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		utils.Logger.Fatalf("Failed to connect to DB: %v", err)
	}
	defer dbPool.Close()

	if err := utils.RunMigrations(dbPool); err != nil {
		utils.Logger.Fatalf("Failed to run migrations: %v", err)
	}

	gws := gateways.NewGateways(cfg.Gateways)

	txnRepo := repository.NewTransactionRepository(dbPool)
	txnService := service.NewTransactionService(txnRepo, gws)

	txnHandler := handlers.NewTransactionHandler(txnService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = utils.NewValidator()

	e.POST("/transactions", txnHandler.CreateTransaction)
	e.POST("/callbacks", txnHandler.HandleCallback)

	go func() {
		mockEcho := echo.New()
		mocks.GatewayAMock(mockEcho)
		mocks.GatewayBMock(mockEcho)
		mockEcho.Start(":8080")
	}()

	go func() {
		if err := e.Start(cfg.Server.Port); err != nil && err != http.ErrServerClosed {
			utils.Logger.Fatalf("Shutting down the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		utils.Logger.Fatal(err)
	}
}
