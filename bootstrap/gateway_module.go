package bootstrap

import (
	"fmt"
	"time"

	gatewayappmax "github.com/appmaxbrasil/appstore-backend-example/app/gateway/appmax"
	gatewaycontracts "github.com/appmaxbrasil/appstore-backend-example/app/gateway/appmax/contracts"
)

type GatewayModule struct {
	AppmaxGateway gatewaycontracts.Gateway
}

func NewGatewayModule(cfg AppmaxConfig) (*GatewayModule, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	client, err := gatewayappmax.NewClientWithOptions(cfg.AuthURL, cfg.APIURL, gatewayappmax.ClientOptions{
		RetryMax:      3,
		RetryWait:     5 * time.Second,
		RetryStatuses: []int{502, 503, 504},
	})
	if err != nil {
		return nil, fmt.Errorf("new gateway module: %w", err)
	}

	return &GatewayModule{AppmaxGateway: client}, nil
}
