package oneplus

import (
	"bytes"
	"encoding/json"
	"fmt"
	shellymodels "github.com/darkmatus/shelly/clients/shellyModels"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

// Shelly1Plus struct represents the state.
type Shelly1Plus struct {
	AuthKey    string
	BaseURL    string
	DeviceID   string
	HTTPClient *http.Client
}

// NewClient Create a new clients with the given params
func NewClient(authKey string, baseURL string, deviceID string) *Shelly1Plus {
	return &Shelly1Plus{
		AuthKey:    authKey,
		BaseURL:    baseURL,
		DeviceID:   deviceID,
		HTTPClient: &http.Client{},
	}
}

// GetDeviceStatus returns the status of the device
func (s *Shelly1Plus) GetDeviceStatus() (shellymodels.Status, error) {
	status := shellymodels.Status{}
	url := fmt.Sprintf("%s/device/status?id=%s&auth_key=%s", s.BaseURL, s.DeviceID, s.AuthKey)

	response, err := s.HTTPClient.Get(url)
	if err != nil {
		return status, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return status, err
	}

	err = json.Unmarshal(body, &status)
	if err != nil {
		return status, err
	}

	return status, nil
}

// On switch the relay on
func (s *Shelly1Plus) On(channelID int) (bool, error) {
	uri := fmt.Sprintf("%s/device/relay/control", s.BaseURL)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// add form data
	_ = writer.WriteField("id", s.DeviceID)
	_ = writer.WriteField("channel", strconv.Itoa(channelID))
	_ = writer.WriteField("turn", "on")
	_ = writer.WriteField("auth_key", s.AuthKey)

	contentType := writer.FormDataContentType()
	_ = writer.Close()

	// create POST-request
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		fmt.Println("error while creating http-request:", err)
		return false, err
	}
	req.Header.Set("Content-Type", contentType)

	// execute HTTP-Request
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		fmt.Println("error while sending http-request:", err)
		return false, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var responseBody = SwitchResponse{}

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		return false, err
	}
	return responseBody.Isok, nil
}

// Off switch the relay off
func (s *Shelly1Plus) Off(channelID int) (bool, error) {

	uri := fmt.Sprintf("%s/device/relay/control", s.BaseURL)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// add form data
	_ = writer.WriteField("id", s.DeviceID)
	_ = writer.WriteField("channel", strconv.Itoa(channelID))
	_ = writer.WriteField("turn", "off")
	_ = writer.WriteField("auth_key", s.AuthKey)

	contentType := writer.FormDataContentType()
	_ = writer.Close()

	// POST-Request erstellen
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		fmt.Println("error while creating http-requestHTTP-Requests:", err)
		return false, err
	}
	req.Header.Set("Content-Type", contentType)

	// execute HTTP-Request
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		fmt.Println("error while sending http-request:", err)
		return false, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var responseBody = SwitchResponse{}

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		return false, err
	}
	return responseBody.Isok, nil
}

// Toggle toggles the current state of the switch
func (s *Shelly1Plus) Toggle(channelID int) (bool, error) {
	uri := fmt.Sprintf("%s/device/relay/control", s.BaseURL)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// add form data
	_ = writer.WriteField("id", s.DeviceID)
	_ = writer.WriteField("channel", strconv.Itoa(channelID))
	_ = writer.WriteField("turn", "toggle")
	_ = writer.WriteField("auth_key", s.AuthKey)

	contentType := writer.FormDataContentType()
	_ = writer.Close()

	// create POST-Request
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		fmt.Println("error while creating http-requestHTTP-Requests:", err)
		return false, err
	}
	req.Header.Set("Content-Type", contentType)

	// execute HTTP-Request
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		fmt.Println("error while sending http-request:", err)
		return false, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var responseBody = SwitchResponse{}

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		return false, err
	}
	return responseBody.Isok, nil
}
