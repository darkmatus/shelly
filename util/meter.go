package util

// Meter provides total active power in W
type Meter interface {
	CurrentPower() (float64, error)
}
