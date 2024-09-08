package gateway_a

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"exinity-golang-assessment/internal/domain"
)

type GatewayA struct {
	name        string
	apiEndpoint string
	apiTimeout  time.Duration
	username    string
	password    string
	client      *http.Client
}

type gatewayARequest struct {
	TransactionID string  `json:"transaction_id"`
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
}

type gatewayAResponse struct {
	ExternalTxnID string `json:"external_txn_id"`
	Status        string `json:"status"`
}

func NewGatewayA(name, apiEndpoint, username, password string, apiTimeout time.Duration) *GatewayA {
	return &GatewayA{
		name:        name,
		apiEndpoint: apiEndpoint,
		apiTimeout:  apiTimeout,
		username:    username,
		password:    password,
		client:      &http.Client{},
	}
}

func (g *GatewayA) Name() string {
	return g.name
}

func (g *GatewayA) Deposit(ctx context.Context, txn *domain.Transaction) (string, error) {
	return g.Request(ctx, txn)
}

func (g *GatewayA) Withdraw(ctx context.Context, txn *domain.Transaction) (string, error) {
	return g.Request(ctx, txn)
}

func (g *GatewayA) Request(ctx context.Context, txn *domain.Transaction) (string, error) {
	reqPayload := gatewayARequest{
		TransactionID: txn.ID,
		Type:          string(txn.Type),
		Amount:        txn.Amount,
		Currency:      txn.Currency,
	}
	data, err := json.Marshal(reqPayload)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/withdraw", g.apiEndpoint)
	if txn.Type == domain.Deposit {
		url = fmt.Sprintf("%s/deposit", g.apiEndpoint)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, g.apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(g.username, g.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		var respPayload gatewayAResponse
		if err := json.NewDecoder(resp.Body).Decode(&respPayload); err != nil {
			return "", err
		}
		return respPayload.ExternalTxnID, nil
	}

	return "", errors.New(resp.Status)
}

func (g *GatewayA) HandleCallback(ctx context.Context, payload map[string]interface{}) (string, domain.TransactionStatus, error) {
	txnID, ok := payload["transaction_id"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid transaction_id")
	}
	statusStr, ok := payload["status"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid status")
	}

	var status domain.TransactionStatus
	switch statusStr {
	case "success":
		status = domain.Completed
	case "failed":
		status = domain.Failed
	default:
		status = domain.Pending
	}

	return txnID, status, nil
}
