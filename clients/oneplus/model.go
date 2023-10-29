package oneplus

type SwitchResponse struct {
	Isok bool       `json:"isok"`
	Data SwitchData `json:"data"`
}
type SwitchData struct {
	DeviceID string `json:"device_id"`
}
