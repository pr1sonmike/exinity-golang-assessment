package gateways

import (
	"exinity-golang-assessment/config"
	"exinity-golang-assessment/internal/gateways/gateway_a"
	"exinity-golang-assessment/internal/gateways/gateway_b"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGateways_Success(t *testing.T) {
	cfgs := []config.GatewayConfig{
		{
			Name:        "GatewayA",
			APIEndpoint: "https://api.gatewaya.com",
			Auth: struct {
				Username string `mapstructure:"username"`
				Password string `mapstructure:"password"`
			}{
				Username: "userA",
				Password: "passA",
			},
			APITimeout: 30,
		},
		{
			Name:        "GatewayB",
			APIEndpoint: "https://api.gatewayb.com",
			Auth: struct {
				Username string `mapstructure:"username"`
				Password string `mapstructure:"password"`
			}{
				Username: "userA",
				Password: "passA",
			},
			APITimeout: 30,
		},
	}

	gateways := NewGateways(cfgs)

	assert.Len(t, gateways, 2)
	assert.IsType(t, &gateway_a.GatewayA{}, gateways[0])
	assert.IsType(t, &gateway_b.GatewayB{}, gateways[1])
}

func TestNewGateways_UnknownGateway(t *testing.T) {
	cfgs := []config.GatewayConfig{
		{
			Name:        "UnknownGateway",
			APIEndpoint: "https://api.unknowngateway.com",
			Auth: struct {
				Username string `mapstructure:"username"`
				Password string `mapstructure:"password"`
			}{
				Username: "userA",
				Password: "passA",
			},
			APITimeout: 30,
		},
	}

	assert.Panics(t, func() {
		NewGateways(cfgs)
	})
}
