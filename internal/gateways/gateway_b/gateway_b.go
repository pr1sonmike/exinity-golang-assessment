package gateway_b

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"time"

	"exinity-golang-assessment/internal/domain"
)

type GatewayB struct {
	name        string
	apiEndpoint string
	apiTimeout  time.Duration
	username    string
	password    string
	client      *http.Client
}

type gatewayBRequest struct {
	XMLName       xml.Name `xml:"TransactionRequest"`
	TransactionID string   `xml:"TransactionID"`
	Type          string   `xml:"Type"`
	Amount        float64  `xml:"Amount"`
	Currency      string   `xml:"Currency"`
}

type gatewayBResponse struct {
	XMLName       xml.Name `xml:"TransactionResponse"`
	ExternalTxnID string   `xml:"ExternalTxnID"`
	Status        string   `xml:"Status"`
}

func NewGatewayB(name, apiEndpoint, username, password string, apiTimeout time.Duration) *GatewayB {
	return &GatewayB{
		name:        name,
		apiEndpoint: apiEndpoint,
		apiTimeout:  apiTimeout,
		username:    username,
		password:    password,
		client:      &http.Client{},
	}
}

func (g *GatewayB) Name() string {
	return g.name
}

func (g *GatewayB) Deposit(ctx context.Context, txn *domain.Transaction) (string, error) {
	return g.sendTransaction(ctx, txn, "Deposit")
}

func (g *GatewayB) Withdraw(ctx context.Context, txn *domain.Transaction) (string, error) {
	return g.sendTransaction(ctx, txn, "Withdraw")
}

func (g *GatewayB) sendTransaction(ctx context.Context, txn *domain.Transaction, txnType string) (string, error) {
	reqPayload := gatewayBRequest{
		TransactionID: txn.ID,
		Type:          txnType,
		Amount:        txn.Amount,
		Currency:      txn.Currency,
	}
	data, err := xml.Marshal(reqPayload)
	if err != nil {
		return "", err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, g.apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, "POST", g.apiEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(g.username, g.password)
	req.Header.Set("Content-Type", "application/xml")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		var respPayload gatewayBResponse
		if err := xml.NewDecoder(resp.Body).Decode(&respPayload); err != nil {
			return "", err
		}

		return respPayload.ExternalTxnID, nil
	}

	return "", errors.New(resp.Status)
}

func (g *GatewayB) HandleCallback(ctx context.Context, payload map[string]interface{}) (string, domain.TransactionStatus, error) {
	txnID, ok := payload["TransactionID"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid TransactionID")
	}
	statusStr, ok := payload["Status"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid Status")
	}

	var status domain.TransactionStatus
	switch statusStr {
	case "Success":
		status = domain.Completed
	case "Failed":
		status = domain.Failed
	default:
		status = domain.Pending
	}

	return txnID, status, nil
}
