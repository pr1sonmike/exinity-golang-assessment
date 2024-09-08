package services

import (
	"errors"
	"exinity-golang-assessment/internal/gateways"
	"exinity-golang-assessment/internal/models"

	"gorm.io/gorm"
)

type TransactionService struct {
	DB       *gorm.DB
	GatewayA *gateways.GatewayA
	GatewayB *gateways.GatewayB
}

type DepositRequest struct {
	Amount float64 `json:"amount"`
}

type WithdrawRequest struct {
	Amount float64 `json:"amount"`
}

type CallbackRequest struct {
	Status string `json:"status"`
}

func NewTransactionService(db *gorm.DB, gatewayA *gateways.GatewayA, gatewayB *gateways.GatewayB) *TransactionService {
	return &TransactionService{DB: db, GatewayA: gatewayA, GatewayB: gatewayB}
}

func (s *TransactionService) Deposit(req *DepositRequest) (*models.Transaction, error) {
	var status string
	var err error
	// Assuming Gateway A for deposits
	status, err = s.GatewayA.Deposit(req.Amount)

	if err != nil {
		return nil, err
	}

	tx := &models.Transaction{Amount: req.Amount, Gateway: "GatewayA", Status: status}
	s.DB.Create(tx)
	return tx, nil
}

func (s *TransactionService) Withdraw(req *WithdrawRequest) (*models.Transaction, error) {
	// Similar logic for withdraw
	return nil, errors.New("Not implemented")
}

func (s *TransactionService) HandleCallback(req *CallbackRequest) error {
	// Update transaction based on callback
	return nil
}
