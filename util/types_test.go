package util

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test Gen2StatusResponse response
func TestUnmarshalGen2StatusResponse(t *testing.T) {
	// Shelly Pro 1PM channel 0 (1)
	var res Gen2StatusResponse

	jsonstr := `{"ble":{},"cloud":{"connected":true},"eth":{"ip":null},"input:0":{"id":0,"state":false},"input:1":{"id":1,"state":false},"mqtt":{"connected":false},"switch:0":{"id":0, "source":"HTTP", "output":false, "apower":47.11, "voltage":232.0, "current":0.000, "pf":0.00, "aenergy":{"total":5.125,"by_minute":[0.000,0.000,0.000],"minute_ts":1675718520},"temperature":{"tC":25.3, "tF":77.5}},"sys":{"mac":"30C6F78BB4D8","restart_required":false,"time":"22:22","unixtime":1675718522,"uptime":45070,"ram_size":234204,"ram_free":137716,"fs_size":524288,"fs_free":172032,"cfg_rev":13,"kvs_rev":1,"schedule_rev":0,"webhook_rev":0,"available_updates":{"beta":{"version":"0.13.0-beta3"}}},"wifi":{"sta_ip":"192.168.178.64","status":"got ip","ssid":"***","rssi":-62},"ws":{"connected":false}}`
	assert.NoError(t, json.Unmarshal([]byte(jsonstr), &res))

	assert.Equal(t, 5.125, res.Switch0.Aenergy.Total)
	assert.Equal(t, 47.11, res.Switch0.Apower)
}
