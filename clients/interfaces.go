package clients

import (
	shellymodels "github.com/darkmatus/shelly/clients/shellyModels"
)

type ShellyInterface interface {
	On(int) (bool, error)
	Off(int) (bool, error)
	GetDeviceStatus() (shellymodels.Status, error)
	Toggle(int) (bool, error)
}
