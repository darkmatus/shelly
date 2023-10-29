package oneplus

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSwitchResponseMarshalling(t *testing.T) {
	// create struct for test
	switchResponse := SwitchResponse{
		Isok: true,
		Data: SwitchData{
			DeviceID: "12345",
		},
	}

	jsonData, err := json.Marshal(switchResponse)
	assert.NoError(t, err)

	var unmarshalledSwitchResponse SwitchResponse
	err = json.Unmarshal(jsonData, &unmarshalledSwitchResponse)
	assert.NoError(t, err)

	assert.Equal(t, switchResponse.Isok, unmarshalledSwitchResponse.Isok)
	assert.Equal(t, switchResponse.Data.DeviceID, unmarshalledSwitchResponse.Data.DeviceID)
}

func TestSwitchResponseUnmarshalling(t *testing.T) {
	jsonData := []byte(`{"isok":true,"data":{"device_id":"54321"}}`)

	var switchResponse SwitchResponse
	err := json.Unmarshal(jsonData, &switchResponse)
	assert.NoError(t, err)

	expectedSwitchResponse := SwitchResponse{
		Isok: true,
		Data: SwitchData{
			DeviceID: "54321",
		},
	}

	assert.Equal(t, expectedSwitchResponse.Isok, switchResponse.Isok)
	assert.Equal(t, expectedSwitchResponse.Data.DeviceID, switchResponse.Data.DeviceID)
}
