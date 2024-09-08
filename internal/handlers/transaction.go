package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"exinity-golang-assessment/internal/domain"
	"exinity-golang-assessment/internal/service"
	"exinity-golang-assessment/internal/utils"
)

type TransactionHandler struct {
	Service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{Service: service}
}

type CreateTransactionRequest struct {
	Type     domain.TransactionType `json:"type" validate:"required,oneof=deposit withdraw"`
	Amount   float64                `json:"amount" validate:"required,gt=0"`
	Currency string                 `json:"currency" validate:"required,len=3"`
	Gateway  string                 `json:"gateway" validate:"required"`
}

func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	var req CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	utils.Logger.Infof("New Transaction: Type: %s, Amount: %f, Currency: %s, Gateway: %s", req.Type, req.Amount, req.Currency, req.Gateway)

	txn, err := h.Service.CreateTransaction(c.Request().Context(), &domain.Transaction{
		Type:     req.Type,
		Amount:   req.Amount,
		Currency: req.Currency,
		Gateway:  req.Gateway,
	})
	if err != nil {
		utils.Logger.Errorf("Failed to create transaction: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	utils.Logger.Infof("Result Transaction: ID: %s, External ID: %s, Status: %s", txn.ID, txn.ExternalTxnID, txn.Status)

	return c.JSON(http.StatusCreated, txn)
}

func (h *TransactionHandler) HandleCallback(c echo.Context) error {
	var payload map[string]interface{}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid callback payload"})
	}

	utils.Logger.Infof("Handling callback: %+v", payload)

	err := h.Service.HandleCallback(c.Request().Context(), payload)
	if err != nil {
		utils.Logger.Errorf("Failed to handle callback: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "success"})
}
