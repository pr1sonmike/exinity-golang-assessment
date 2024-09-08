package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"exinity-golang-assessment/internal/domain"
)

type TransactionRepository interface {
	Create(ctx context.Context, txn *domain.Transaction) error
	UpdateStatus(ctx context.Context, id string, status domain.TransactionStatus) error
	GetByID(ctx context.Context, id string) (*domain.Transaction, error)
}

type transactionRepo struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(ctx context.Context, txn *domain.Transaction) error {
	txn.ID = uuid.New().String()
	txn.Status = domain.Pending
	txn.CreatedAt = time.Now()
	txn.UpdatedAt = time.Now()

	query := `
        INSERT INTO transactions (id, type, amount, currency, gateway, status, created_at, updated_at, external_txn_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	_, err := r.db.Exec(ctx, query,
		txn.ID, txn.Type, txn.Amount, txn.Currency,
		txn.Gateway, txn.Status, txn.CreatedAt, txn.UpdatedAt,
		txn.ExternalTxnID,
	)
	return err
}

func (r *transactionRepo) UpdateStatus(ctx context.Context, id string, status domain.TransactionStatus) error {
	query := `
        UPDATE transactions
        SET status = $1, updated_at = $2
        WHERE id = $3
    `
	cmdTag, err := r.db.Exec(ctx, query, status, time.Now(), id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() != 1 {
		return errors.New("transaction not found")
	}
	return nil
}

func (r *transactionRepo) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	var txn domain.Transaction
	query := `
        SELECT id, type, amount, currency, gateway, status, created_at, updated_at, external_txn_id
        FROM transactions
        WHERE id = $1
    `
	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(
		&txn.ID, &txn.Type, &txn.Amount, &txn.Currency,
		&txn.Gateway, &txn.Status, &txn.CreatedAt,
		&txn.UpdatedAt, &txn.ExternalTxnID,
	)
	if err != nil {
		return nil, err
	}
	return &txn, nil
}
