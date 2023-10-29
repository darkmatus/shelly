package ShellyModels

type Status struct {
	Isok bool `json:"isok"`
	Data Data `json:"data"`
}
type Mqtt struct {
	Connected bool `json:"connected"`
}
type Temperature struct {
	TC float64 `json:"tC"`
	TF float64 `json:"tF"`
}
type Switch0 struct {
	ID          int         `json:"id"`
	Source      string      `json:"source"`
	Output      bool        `json:"output"`
	Temperature Temperature `json:"temperature"`
}
type Ws struct {
	Connected bool `json:"connected"`
}
type VEve0 struct {
	Ev  string `json:"ev"`
	TTL int    `json:"ttl"`
	ID  int    `json:"id"`
}
type Beta struct {
	Version string `json:"version"`
}
type AvailableUpdates struct {
	Beta Beta `json:"beta"`
}
type Sys struct {
	Mac              string           `json:"mac"`
	RestartRequired  bool             `json:"restart_required"`
	Time             string           `json:"time"`
	Unixtime         int              `json:"unixtime"`
	Uptime           int              `json:"uptime"`
	RAMSize          int              `json:"ram_size"`
	RAMFree          int              `json:"ram_free"`
	FsSize           int              `json:"fs_size"`
	FsFree           int              `json:"fs_free"`
	CfgRev           int              `json:"cfg_rev"`
	KvsRev           int              `json:"kvs_rev"`
	ScheduleRev      int              `json:"schedule_rev"`
	WebhookRev       int              `json:"webhook_rev"`
	AvailableUpdates AvailableUpdates `json:"available_updates"`
}
type Cloud struct {
	Connected bool `json:"connected"`
}
type Wifi struct {
	StaIP  string `json:"sta_ip"`
	Status string `json:"status"`
	Ssid   string `json:"ssid"`
	Rssi   int    `json:"rssi"`
}
type Input0 struct {
	ID    int  `json:"id"`
	State bool `json:"state"`
}
type DeviceStatus struct {
	Serial  int     `json:"serial"`
	Mqtt    Mqtt    `json:"mqtt"`
	Switch0 Switch0 `json:"switch:0"`
	ID      string  `json:"id"`
	Ws      Ws      `json:"ws"`
	VEve0   VEve0   `json:"v_eve:0"`
	Sys     Sys     `json:"sys"`
	Ble     []any   `json:"ble"`
	Cloud   Cloud   `json:"cloud"`
	Wifi    Wifi    `json:"wifi"`
	Input0  Input0  `json:"input:0"`
	Updated string  `json:"_updated"`
	Code    string  `json:"code"`
}
type Data struct {
	Online       bool         `json:"online"`
	DeviceStatus DeviceStatus `json:"device_status"`
}
