package shelly

import (
	"errors"
	"fmt"
	"github.com/darkmatus/shelly/connection"
	"github.com/darkmatus/shelly/util"
	"math"
)

type Switch struct {
	*connection.Connection
}

func NewSwitch(conn *connection.Connection) *Switch {
	res := &Switch{
		Connection: conn,
	}

	return res
}

// CurrentPower implements the api.Meter interface
func (sh *Switch) CurrentPower() (float64, error) {
	var power float64

	d := sh.Connection
	switch d.Gen {
	case 0, 1:
		var res util.Gen1StatusResponse
		uri := fmt.Sprintf("%s/status", d.URI)
		if err := d.GetJSON(uri, &res); err != nil {
			return 0, err
		}

		switch {
		case d.Channel < len(res.Meters):
			power = res.Meters[d.Channel].Power
		case d.Channel < len(res.EMeters):
			power = res.EMeters[d.Channel].Power
		default:
			return 0, errors.New("invalid channel, missing power meter")
		}

	default:
		var res util.Gen2StatusResponse
		if err := d.ExecGen2Cmd("Shelly.GetStatus", false, &res); err != nil {
			return 0, err
		}

		switch d.Channel {
		case 1:
			power = res.Switch1.Apower
		case 2:
			power = res.Switch2.Apower
		default:
			power = res.Switch0.Apower
		}
	}

	// Assure positive power response (Gen 1 EM devices can provide negative values)
	return math.Abs(power), nil
}

// Enabled implements the api.Charger interface
func (sh *Switch) Enabled() (bool, error) {
	d := sh.Connection
	switch d.Gen {
	case 0, 1:
		var res util.Gen1SwitchResponse
		uri := fmt.Sprintf("%s/relay/%d", d.URI, d.Channel)
		err := d.GetJSON(uri, &res)
		return res.Ison, err

	default:
		var res util.Gen2SwitchResponse
		err := d.ExecGen2Cmd("Switch.GetStatus", false, &res)
		return res.Output, err
	}
}

// Enable implements the Charger interface
func (sh *Switch) Enable(enable bool) error {
	var err error
	onoff := map[bool]string{true: "on", false: "off"}

	d := sh.Connection
	switch d.Gen {
	case 0, 1:
		var res util.Gen1SwitchResponse
		uri := fmt.Sprintf("%s/relay/%d?turn=%s", d.URI, d.Channel, onoff[enable])
		err = d.GetJSON(uri, &res)

	default:
		var res util.Gen2SwitchResponse
		err = d.ExecGen2Cmd("Switch.Set", enable, &res)
	}

	return err
}
