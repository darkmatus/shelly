package shelly

import (
	"fmt"
	"github.com/darkmatus/shelly/clients"
	"github.com/darkmatus/shelly/clients/oneplus"
)

const (
	Device1Plus = "1plus"
)

// NewShelly creates a Shelly charger from generic config
func NewShelly(deviceType string, authKey string, baseURL string, deviceID string) (clients.ShellyInterface, error) {
	switch deviceType {
	case Device1Plus:
		return oneplus.NewClient(authKey, baseURL, deviceID), nil
	default:
		return nil, fmt.Errorf("given device type '%s' is currently not supported", deviceType)
	}
}
