package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/darkmatus/shelly/util"
)

// Timeout is the default request timeout used by the Helper
var Timeout = 10 * time.Second

// Helper provides utility primitives
type Helper struct {
	*http.Client
}

// NewClient creates http client with default transport
func NewClient(log *util.Logger) *http.Client {
	return &http.Client{
		Timeout:   Timeout,
		Transport: NewTripper(log, util.Default()),
	}
}

// NewHelper creates http helper for simplified PUT GET logic
func NewHelper(log *util.Logger) *Helper {
	return &Helper{
		Client: NewClient(log),
	}
}

// DoBody executes HTTP request and returns the response body
func (r *Helper) DoBody(req *http.Request) ([]byte, error) {
	resp, err := r.Do(req)
	var body []byte
	if err == nil {
		body, err = ReadBody(resp)
	}
	return body, err
}

// GetBody executes HTTP GET request and returns the response body
func (r *Helper) GetBody(url string) ([]byte, error) {
	resp, err := r.Get(url)
	var body []byte
	if err == nil {
		body, err = ReadBody(resp)
	}
	return body, err
}

// decodeJSON reads HTTP response and decodes JSON body if error is nil
func decodeJSON(resp *http.Response, res interface{}) error {
	if err := ResponseError(resp); err != nil {
		_ = json.NewDecoder(resp.Body).Decode(&res)
		return err
	}

	return json.NewDecoder(resp.Body).Decode(&res)
}

// DoJSON executes HTTP request and decodes JSON response.
// It returns a StatusError on response codes other than HTTP 2xx.
func (r *Helper) DoJSON(req *http.Request, res interface{}) error {
	resp, err := r.Do(req)
	if err == nil {
		defer resp.Body.Close()
		err = decodeJSON(resp, &res)
	}
	return err
}

// GetJSON executes HTTP GET request and decodes JSON response.
// It returns a StatusError on response codes other than HTTP 2xx.
func (r *Helper) GetJSON(url string, res interface{}) error {
	req, err := New(http.MethodGet, url, nil, AcceptJSON)
	if err == nil {
		err = r.DoJSON(req, &res)
	}
	return err
}

// DefaultScheme prepends given scheme to uri if not specified
func DefaultScheme(uri, scheme string) string {
	u, err := url.Parse(uri)
	if err != nil {
		if strings.HasSuffix(err.Error(), "first path segment in URL cannot contain colon") {
			return fmt.Sprintf("%s://%s", scheme, uri)
		}
		return uri
	}

	switch {
	case u.String() == "":
		return ""

	case u.Scheme == "":
		// scheme missing
		u.Scheme = scheme

	case u.Opaque != "":
		// host:port format is parsed as scheme:opaque (https://golang.org/pkg/net/url/#URL)
		if u, err = url.Parse(fmt.Sprintf("%s://%s", scheme, uri)); err != nil {
			return uri
		}
	}

	return u.String()
}
