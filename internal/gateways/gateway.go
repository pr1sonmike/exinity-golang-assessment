package gateways

import (
	"context"

	"exinity-golang-assessment/internal/domain"
)

type PaymentGateway interface {
	Name() string
	Deposit(ctx context.Context, txn *domain.Transaction) (string, error)
	Withdraw(ctx context.Context, txn *domain.Transaction) (string, error)
	HandleCallback(ctx context.Context, payload map[string]interface{}) (string, domain.TransactionStatus, error)
}
