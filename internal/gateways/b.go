package gateways

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type GatewayB struct {
	Client *http.Client
}

func NewGatewayB() *GatewayB {
	return &GatewayB{
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (g *GatewayB) Deposit(amount float64) (string, error) {
	soapBody := fmt.Sprintf(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
		<soapenv:Body>
			<DepositRequest>
				<Amount>%f</Amount>
			</DepositRequest>
		</soapenv:Body>
	</soapenv:Envelope>`, amount)

	req, err := http.NewRequest("POST", "https://gatewayB.com/deposit", bytes.NewBuffer([]byte(soapBody)))
	req.Header.Set("Content-Type", "text/xml")
	if err != nil {
		return "", err
	}

	resp, err := g.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Handle SOAP/XML parsing and status extraction here.
	return "Success", nil
}
