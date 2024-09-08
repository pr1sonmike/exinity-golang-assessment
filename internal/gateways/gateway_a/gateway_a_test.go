package gateway_a

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"exinity-golang-assessment/internal/domain"
)

func mockTransaction() *domain.Transaction {
	return &domain.Transaction{
		ID:       "12345",
		Type:     domain.Deposit,
		Amount:   100.50,
		Currency: "USD",
	}
}

func TestGatewayA_Deposit(t *testing.T) {
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
				w.Write([]byte(`{"external_txn_id":"ext12345","status":"success"}`))
			},
			expectedErr: nil,
			expectedID:  "ext12345",
			txn:         mockTransaction(),
		},
		{
			name: "Failed Deposit - Invalid Transaction",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"failed"}`))
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

			gateway := NewGatewayA("GatewayA", server.URL, "username", "password", time.Duration(1)*time.Second)
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

func TestGatewayA_HandleCallback(t *testing.T) {
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
				"transaction_id": "12345",
				"status":         "success",
			},
		},
		{
			name:           "Failed Callback",
			expectedErr:    nil,
			expectedID:     "12345",
			expectedStatus: "failed",
			payload: map[string]interface{}{
				"transaction_id": "12345",
				"status":         "failed",
			},
		},
		{
			name:           "Failed Callback",
			expectedErr:    nil,
			expectedID:     "12345",
			expectedStatus: "pending",
			payload: map[string]interface{}{
				"transaction_id": "12345",
				"status":         "pending",
			},
		},
		{
			name:        "Invalid Transaction",
			expectedErr: errors.New("invalid transaction_id"),
			expectedID:  "",
			payload: map[string]interface{}{
				"status": "success",
			},
		},
		{
			name:        "Invalid Status",
			expectedErr: errors.New("invalid status"),
			expectedID:  "",
			payload: map[string]interface{}{
				"transaction_id": "12345",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gateway := NewGatewayA("GatewayA", "test", "username", "password", 1)
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
