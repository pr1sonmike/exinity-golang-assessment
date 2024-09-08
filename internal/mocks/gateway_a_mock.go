package mocks

import (
	"exinity-golang-assessment/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GatewayAMock(e *echo.Echo) {
	e.POST("/api/deposit", apiHandler)
	e.POST("/api/withdraw", apiHandler)
}

func apiHandler(c echo.Context) error {
	var req map[string]interface{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	utils.Logger.Infof("GatewayA: URL: %s, Request: %v", c.Request().URL, req)
	resp := map[string]string{
		"external_txn_id": "A-" + req["transaction_id"].(string),
		"status":          "pending",
	}
	return c.JSON(http.StatusOK, resp)
}
