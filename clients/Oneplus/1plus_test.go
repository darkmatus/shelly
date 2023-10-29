package Oneplus

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func getTestServer(t *testing.T, path string, method string, responseFile string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the URL and method
		assert.Equal(t, r.URL.Path, path)
		assert.Equal(t, r.Method, method)

		// Simulate a response
		jsonResponse, err := os.ReadFile(responseFile)
		if err != nil {
			t.Fatal(err)
		}

		// Send the response
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}))
	return server
}

func TestGetDeviceStatusOnline(t *testing.T) {

	server := getTestServer(t, "/device/status", "GET", "testdata/getStatusOnline.json")
	defer server.Close()
	//Create an instance of Shelly1Plus with the mock server URL
	shelly := &Shelly1Plus{
		BaseURL:    server.URL,
		DeviceID:   "048961276",
		AuthKey:    "GREATSECRETHASH123",
		HTTPClient: &http.Client{},
	}

	// Call the method to be tested
	status, err := shelly.GetDeviceStatus()

	// Check the result
	assert.NoError(t, err)

	assert.Equal(t, true, status.Data.Online)
}
func TestGetDeviceStatusOffline(t *testing.T) {
	server := getTestServer(t, "/device/status", "GET", "testdata/getStatusOffline.json")
	defer server.Close()

	//Create an instance of Shelly1Plus with the mock server URL
	shelly := &Shelly1Plus{
		BaseURL:    server.URL,
		DeviceID:   "048961276",
		AuthKey:    "GREATSECRETHASH123",
		HTTPClient: &http.Client{},
	}

	// Call the method to be tested
	status, err := shelly.GetDeviceStatus()

	// Check the result
	assert.NoError(t, err)
	// Call the method to be tested
	status, err = shelly.GetDeviceStatus()

	// Check the result
	assert.NoError(t, err)
	fmt.Println(status.Data.Online)
	assert.Equal(t, false, status.Data.Online)
}

func TestOn(t *testing.T) {
	server := getTestServer(t, "/device/relay/control", "POST", "testdata/response-on.json")
	defer server.Close()
	//Create an instance of Shelly1Plus with the mock server URL
	shelly := &Shelly1Plus{
		BaseURL:    server.URL,
		DeviceID:   "048961276",
		AuthKey:    "GREATSECRETHASH123",
		HTTPClient: &http.Client{},
	}

	// Call the method to be tested
	res, err := shelly.On(0)

	assert.NoError(t, err)
	assert.Equal(t, true, res)
}

func TestOff(t *testing.T) {

	server := getTestServer(t, "/device/relay/control", "POST", "testdata/response-off.json")
	defer server.Close()
	//Create an instance of Shelly1Plus with the mock server URL
	shelly := &Shelly1Plus{
		BaseURL:    server.URL,
		DeviceID:   "048961276",
		AuthKey:    "GREATSECRETHASH123",
		HTTPClient: &http.Client{},
	}

	// Call the method to be tested
	res, err := shelly.Off(0)

	assert.NoError(t, err)
	assert.Equal(t, true, res)
}

func TestToggle(t *testing.T) {

	server := getTestServer(t, "/device/relay/control", "POST", "testdata/response-on.json")
	defer server.Close()
	//Create an instance of Shelly1Plus with the mock server URL
	shelly := &Shelly1Plus{
		BaseURL:    server.URL,
		DeviceID:   "048961276",
		AuthKey:    "GREATSECRETHASH123",
		HTTPClient: &http.Client{},
	}

	// Call the method to be tested
	res, err := shelly.Toggle(0)

	assert.NoError(t, err)
	assert.Equal(t, true, res)
}
