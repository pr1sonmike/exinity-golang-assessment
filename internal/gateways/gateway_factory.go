package gateways

import (
	"exinity-golang-assessment/config"
	"exinity-golang-assessment/internal/gateways/gateway_a"
	"exinity-golang-assessment/internal/gateways/gateway_b"
	"exinity-golang-assessment/internal/utils"
)

func NewGateways(cfgs []config.GatewayConfig) []PaymentGateway {
	var gateways []PaymentGateway
	for _, cfg := range cfgs {
		switch cfg.Name {
		case "GatewayA":
			gateways = append(gateways, gateway_a.NewGatewayA(cfg.Name, cfg.APIEndpoint, cfg.Auth.Username, cfg.Auth.Password, cfg.APITimeout))
		case "GatewayB":
			gateways = append(gateways, gateway_b.NewGatewayB(cfg.Name, cfg.APIEndpoint, cfg.Auth.Username, cfg.Auth.Password, cfg.APITimeout))
		default:
			utils.Logger.Fatal("Unknown gateway")
		}
	}
	return gateways
}
