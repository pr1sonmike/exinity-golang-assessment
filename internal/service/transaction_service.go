package service

import (
	"context"
	"errors"
	"exinity-golang-assessment/internal/utils"
	"strings"
	"time"

	"exinity-golang-assessment/internal/domain"
	"exinity-golang-assessment/internal/gateways"
	"exinity-golang-assessment/internal/repository"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, txn *domain.Transaction) (*domain.Transaction, error)
	HandleCallback(ctx context.Context, payload map[string]interface{}) error
}

type transactionService struct {
	repo     repository.TransactionRepository
	gateways map[string]gateways.PaymentGateway
}

func NewTransactionService(repo repository.TransactionRepository, gws []gateways.PaymentGateway) TransactionService {
	gwMap := make(map[string]gateways.PaymentGateway)
	for _, gw := range gws {
		gwMap[strings.ToLower(gw.Name())] = gw
	}
	return &transactionService{
		repo:     repo,
		gateways: gwMap,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, txn *domain.Transaction) (*domain.Transaction, error) {
	gw, exists := s.gateways[strings.ToLower(txn.Gateway)]
	if !exists {
		return nil, errors.New("unsupported gateway")
	}

	if err := s.repo.Create(ctx, txn); err != nil {
		return nil, err
	}

	var externalTxnID string
	var err error

	maxRetries := 3
	initialBackoff := 500 * time.Millisecond

	err = utils.Retry(ctx, maxRetries, initialBackoff, func() error {
		switch txn.Type {
		case domain.Deposit:
			externalTxnID, err = gw.Deposit(ctx, txn)
		case domain.Withdraw:
			externalTxnID, err = gw.Withdraw(ctx, txn)
		default:
			return errors.New("invalid transaction type")
		}
		return err
	})

	if err != nil {
		s.repo.UpdateStatus(ctx, txn.ID, domain.Failed)
		return nil, err
	}

	txn.ExternalTxnID = externalTxnID
	s.repo.UpdateStatus(ctx, txn.ID, domain.Processing)

	return txn, nil
}

func (s *transactionService) HandleCallback(ctx context.Context, payload map[string]interface{}) error {
	gatewayName, ok := payload["gateway"].(string)
	if !ok {
		return errors.New("gateway not specified in callback")
	}

	gw, exists := s.gateways[strings.ToLower(gatewayName)]
	if !exists {
		return errors.New("unsupported gateway in callback")
	}

	txnID, status, err := gw.HandleCallback(ctx, payload)
	if err != nil {
		return err
	}

	return s.repo.UpdateStatus(ctx, txnID, status)
}
