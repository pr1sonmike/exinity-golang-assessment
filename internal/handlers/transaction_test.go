package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"exinity-golang-assessment/internal/domain"
	"exinity-golang-assessment/internal/utils"
)

// MockTransactionService is a mock implementation of the TransactionService interface
type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) CreateTransaction(ctx context.Context, txn *domain.Transaction) (*domain.Transaction, error) {
	args := m.Called(ctx, txn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionService) HandleCallback(ctx context.Context, payload map[string]interface{}) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func TestCreateTransaction(t *testing.T) {
	e := echo.New()
	e.Validator = utils.NewValidator()
	utils.InitLogger()
	mockService := new(MockTransactionService)
	handler := NewTransactionHandler(mockService)

	t.Run("valid request", func(t *testing.T) {
		reqBody := `{"type":"deposit","amount":100.0,"currency":"USD","gateway":"A"}`
		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader([]byte(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockTxn := &domain.Transaction{
			ID:            "txn_123",
			Type:          domain.Deposit,
			Amount:        100.0,
			Currency:      "USD",
			Gateway:       "A",
			ExternalTxnID: "ext_123",
			Status:        domain.Pending,
		}
		mockService.
			On("CreateTransaction", mock.Anything, mock.AnythingOfType("*domain.Transaction")).
			Return(mockTxn, nil).
			Once()

		if assert.NoError(t, handler.CreateTransaction(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			var resp domain.Transaction
			if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp)) {
				assert.Equal(t, "txn_123", resp.ID)
				assert.Equal(t, "ext_123", resp.ExternalTxnID)
				assert.Equal(t, domain.Pending, resp.Status)
			}
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		reqBody := `{"type":"invalid","amount":-100.0,"currency":"US","gateway":""}`
		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader([]byte(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, handler.CreateTransaction(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		reqBody := `{"type":"deposit","amount":100.0,"currency":"USD","gateway":"A"}`
		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader([]byte(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.
			On("CreateTransaction", mock.Anything, mock.AnythingOfType("*domain.Transaction")).
			Return(nil, assert.AnError).
			Once()

		if assert.NoError(t, handler.CreateTransaction(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestHandleCallback(t *testing.T) {
	e := echo.New()
	e.Validator = utils.NewValidator()
	utils.InitLogger()
	mockService := new(MockTransactionService)
	handler := NewTransactionHandler(mockService)

	t.Run("valid callback", func(t *testing.T) {
		payload := map[string]interface{}{
			"transaction_id": "txn_123",
			"status":         "success",
		}
		reqBody, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/callback", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.
			On("HandleCallback", mock.Anything, payload).
			Return(nil).
			Once()

		if assert.NoError(t, handler.HandleCallback(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("invalid callback payload", func(t *testing.T) {
		reqBody := `{"transaction_id":123, status:true}`
		req := httptest.NewRequest(http.MethodPost, "/callback", bytes.NewReader([]byte(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, handler.HandleCallback(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		payload := map[string]interface{}{
			"transaction_id": "txn_123",
			"status":         "success",
		}
		reqBody, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/callback", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("HandleCallback", mock.Anything, payload).Return(assert.AnError)

		if assert.NoError(t, handler.HandleCallback(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}
