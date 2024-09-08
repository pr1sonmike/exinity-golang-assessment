package domain

import (
	"time"
)

type TransactionType string
type TransactionStatus string

const (
	Deposit  TransactionType = "deposit"
	Withdraw TransactionType = "withdraw"

	Pending    TransactionStatus = "pending"
	Completed  TransactionStatus = "completed"
	Failed     TransactionStatus = "failed"
	Processing TransactionStatus = "processing"
)

type Transaction struct {
	ID            string            `json:"id"`
	Type          TransactionType   `json:"type"`
	Amount        float64           `json:"amount"`
	Currency      string            `json:"currency"`
	Gateway       string            `json:"gateway"`
	Status        TransactionStatus `json:"status"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	ExternalTxnID string            `json:"external_txn_id"`
}
