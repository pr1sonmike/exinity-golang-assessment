package mocks

import (
	"encoding/xml"
	"net/http"

	"github.com/labstack/echo/v4"

	"exinity-golang-assessment/internal/utils"
)

type TransactionRequest struct {
	XMLName       xml.Name `xml:"TransactionRequest"`
	TransactionID string   `xml:"TransactionID"`
	Type          string   `xml:"Type"`
	Amount        float64  `xml:"Amount"`
	Currency      string   `xml:"Currency"`
}

type TransactionResponse struct {
	XMLName       xml.Name `xml:"TransactionResponse"`
	ExternalTxnID string   `xml:"ExternalTxnID"`
	Status        string   `xml:"Status"`
}

func GatewayBMock(e *echo.Echo) {
	e.POST("/soap", func(c echo.Context) error {
		var req TransactionRequest
		if err := xml.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return c.String(http.StatusBadRequest, "Invalid XML")
		}
		utils.Logger.Infof("GatewayB: URL: %s, Request: %v", c.Request().URL, req)
		resp := TransactionResponse{
			ExternalTxnID: "B-" + req.TransactionID,
			Status:        "Pending",
		}
		c.Response().Header().Set("Content-Type", "application/xml")
		return c.XML(http.StatusOK, resp)
	})
}
