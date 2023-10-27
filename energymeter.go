package shelly

import (
	"github.com/darkmatus/shelly/connection"
	"github.com/darkmatus/shelly/util"
)

type EnergyMeter struct {
	*connection.Connection
}

// PhaseCurrents provides per-phase current A
type PhaseCurrents interface {
	Currents() (float64, float64, float64, error)
}

// PhaseVoltages provides per-phase voltage V
type PhaseVoltages interface {
	Voltages() (float64, float64, float64, error)
}

// PhasePowers provides signed per-phase power W
type PhasePowers interface {
	Powers() (float64, float64, float64, error)
}

func NewEnergyMeter(conn *connection.Connection) *EnergyMeter {
	res := &EnergyMeter{
		Connection: conn,
	}

	return res
}

// CurrentPower implements the api.Meter interface
func (sh *EnergyMeter) CurrentPower() (float64, error) {
	var res util.Gen2EmStatusResponse
	if err := sh.Connection.ExecGen2Cmd("EM.GetStatus", false, &res); err != nil {
		return 0, err
	}

	return res.TotalPower, nil
}

// TotalEnergy implements the api.Meter interface
func (sh *EnergyMeter) TotalEnergy() (float64, error) {
	var res util.Gen2EmDataStatusResponse
	if err := sh.Connection.ExecGen2Cmd("EMData.GetStatus", false, &res); err != nil {
		return 0, err
	}

	return res.TotalEnergy / 1000, nil
}

var _ PhaseCurrents = (*EnergyMeter)(nil)

// Currents implements the api.PhaseCurrents interface
func (sh *EnergyMeter) Currents() (float64, float64, float64, error) {
	var res util.Gen2EmStatusResponse
	if err := sh.Connection.ExecGen2Cmd("EM.GetStatus", false, &res); err != nil {
		return 0, 0, 0, err
	}

	return res.CurrentA, res.CurrentB, res.CurrentC, nil
}

var _ PhaseVoltages = (*EnergyMeter)(nil)

// Voltages implements the api.PhaseVoltages interface
func (sh *EnergyMeter) Voltages() (float64, float64, float64, error) {
	var res util.Gen2EmStatusResponse
	if err := sh.Connection.ExecGen2Cmd("EM.GetStatus", false, &res); err != nil {
		return 0, 0, 0, err
	}

	return res.VoltageA, res.VoltageB, res.VoltageC, nil
}

var _ PhasePowers = (*EnergyMeter)(nil)

// Powers implements the api.PhasePowers interface
func (sh *EnergyMeter) Powers() (float64, float64, float64, error) {
	var res util.Gen2EmStatusResponse
	if err := sh.Connection.ExecGen2Cmd("EM.GetStatus", false, &res); err != nil {
		return 0, 0, 0, err
	}

	return res.PowerA, res.PowerB, res.PowerC, nil
}
