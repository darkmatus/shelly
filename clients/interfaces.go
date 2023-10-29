package clients

import (
	ShellyModels "github.com/darkmatus/shelly/clients/shellyModels"
)

type ShellyInterface interface {
	On(int) (bool, error)
	Off(int) (bool, error)
	GetDeviceStatus() (ShellyModels.Status, error)
	Toggle(int) (bool, error)
}
