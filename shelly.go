package shelly

import (
	"github.com/darkmatus/shelly/connection"
	"github.com/darkmatus/shelly/util"
)

// NewShelly creates a Shelly charger from generic config
func NewShelly(URI string, User string, Password string, Channel int) (util.Meter, error) {

	conn, err := connection.NewConnection(URI, User, Password, Channel)
	if err != nil {
		return nil, err
	}

	return NewSwitch(conn), nil
}

// NewShellyEnergyMeter creates a Shelly from given data
func NewShellyEnergyMeter(URI string, User string, Password string, Channel int) (util.Meter, error) {
	conn, err := connection.NewConnection(URI, User, Password, Channel)
	if err != nil {
		return nil, err
	}

	return NewEnergyMeter(conn), nil
}
