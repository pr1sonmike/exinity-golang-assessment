package gateways

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type GatewayA struct {
	Client *http.Client
}

func NewGatewayA() *GatewayA {
	return &GatewayA{
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (g *GatewayA) Deposit(amount float64) (string, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"amount": amount,
	})

	resp, err := g.Client.Post("https://gatewayA.com/deposit", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result["status"].(string), nil
}

func (g *GatewayA) Withdraw(amount float64) (string, error) {
	// Similar to Deposit
	return "Success", nil
}
