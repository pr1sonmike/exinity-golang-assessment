package gateway_b

import (
	"context"
	"errors"
	"exinity-golang-assessment/internal/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func mockTransaction() *domain.Transaction {
	return &domain.Transaction{
		ID:       "54321",
		Type:     domain.Deposit,
		Amount:   200.50,
		Currency: "EUR",
	}
}

func TestGatewayB_Deposit(t *testing.T) {
	tests := []struct {
		name        string
		serverFunc  http.HandlerFunc
		expectedErr error
		expectedID  string
		txn         *domain.Transaction
	}{
		{
			name: "Successful Deposit",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<TransactionResponse><ExternalTxnID>ext54321</ExternalTxnID><Status>success</Status></TransactionResponse>`))
			},
			expectedErr: nil,
			expectedID:  "ext54321",
			txn:         mockTransaction(),
		},
		{
			name: "Failed Deposit - Invalid Transaction",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`<TransactionResponse><Status>failed</Status></TransactionResponse>`))
			},
			expectedErr: errors.New("400 Bad Request"),
			expectedID:  "",
			txn:         mockTransaction(),
		},
		{
			name: "Server Error",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedErr: errors.New("500 Internal Server Error"),
			expectedID:  "",
			txn:         mockTransaction(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.serverFunc)
			defer server.Close()

			gateway := NewGatewayB("GatewayB", server.URL, "username", "password", time.Duration(1)*time.Second)
			externalID, err := gateway.Deposit(context.Background(), tt.txn)

			if tt.expectedErr != nil {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErr.Error()) {
					t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
				}
			} else if externalID != tt.expectedID {
				t.Errorf("expected externalID: %v, got: %v", tt.expectedID, externalID)
			}
		})
	}
}

func TestGatewayB_HandleCallback(t *testing.T) {
	tests := []struct {
		name           string
		expectedErr    error
		expectedID     string
		expectedStatus domain.TransactionStatus
		payload        map[string]interface{}
	}{
		{
			name:           "Successful Callback",
			expectedErr:    nil,
			expectedID:     "12345",
			expectedStatus: "completed",
			payload: map[string]interface{}{
				"TransactionID": "12345",
				"Status":        "Success",
			},
		},
		{
			name:           "Failed Callback",
			expectedErr:    nil,
			expectedID:     "12345",
			expectedStatus: "failed",
			payload: map[string]interface{}{
				"TransactionID": "12345",
				"Status":        "Failed",
			},
		},
		{
			name:           "Failed Callback",
			expectedErr:    nil,
			expectedID:     "12345",
			expectedStatus: "pending",
			payload: map[string]interface{}{
				"TransactionID": "12345",
				"Status":        "Pending",
			},
		},
		{
			name:        "Invalid Transaction",
			expectedErr: errors.New("invalid TransactionID"),
			expectedID:  "",
			payload: map[string]interface{}{
				"status": "success",
			},
		},
		{
			name:        "Invalid Status",
			expectedErr: errors.New("invalid Status"),
			expectedID:  "",
			payload: map[string]interface{}{
				"TransactionID": "12345",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gateway := NewGatewayB("GatewayA", "test", "username", "password", 1)
			externalID, status, err := gateway.HandleCallback(context.Background(), tt.payload)
			if tt.expectedErr != nil {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErr.Error()) {
					t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
				}
			} else if externalID != tt.expectedID {
				t.Errorf("expected externalID: %v, got: %v", tt.expectedID, externalID)
			}

			if status != tt.expectedStatus {
				t.Errorf("expected status: %v, got: %v", tt.expectedStatus, status)
			}
		})
	}
}
