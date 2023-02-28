package quickbooks

import (
	"fmt"
)

type PlatformClient struct {
	// The set of quickbooks APIs
	discoveryAPI *DiscoveryAPI
	// The client Id
	clientId string
	// The client Secret
	clientSecret string
	// The minor version of the QB API
	minorVersion string
}

// NewClient initializes a new QuickBooks client for interacting with their Online API
func NewPlatformClient(clientId string, clientSecret string, isProduction bool, minorVersion string) (c *PlatformClient, err error) {
	if minorVersion == "" {
		minorVersion = "65"
	}

	client := PlatformClient{
		clientId:     clientId,
		clientSecret: clientSecret,
		minorVersion: minorVersion,
	}

	if isProduction {
		client.discoveryAPI, err = CallDiscoveryAPI(DiscoveryProductionEndpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to obtain discovery endpoint: %v", err)
		}
	} else {
		client.discoveryAPI, err = CallDiscoveryAPI(DiscoverySandboxEndpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to obtain discovery endpoint: %v", err)
		}
	}

	return &client, nil
}
