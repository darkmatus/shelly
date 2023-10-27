package connection

import (
	"errors"
	"fmt"
	"github.com/darkmatus/shelly/request"
	"github.com/darkmatus/shelly/util"
	"github.com/jpfielding/go-http-digest/pkg/digest"
	"net/http"
	"strings"
)

// Connection is the Shelly connection
type Connection struct {
	*request.Helper
	URI        string
	Channel    int
	Gen        int    // Shelly api generation
	Devicetype string // Shelly device type
}

// NewConnection creates a new Shelly device connection.
func NewConnection(uri, user, password string, channel int) (*Connection, error) {
	if uri == "" {
		return nil, errors.New("missing uri")
	}

	for _, suffix := range []string{"/", "/rcp", "/shelly"} {
		uri = strings.TrimSuffix(uri, suffix)
	}

	log := util.NewLogger("shelly")
	client := request.NewHelper(log)

	// Shelly Gen1 and Gen2 families expose the /shelly endpoint
	var resp util.DeviceInfo
	if err := client.GetJSON(fmt.Sprintf("%s/shelly", request.DefaultScheme(uri, "http")), &resp); err != nil {
		return nil, err
	}

	conn := &Connection{
		Helper:     client,
		Channel:    channel,
		Gen:        resp.Gen,
		Devicetype: resp.Type,
	}

	conn.Client.Transport = request.NewTripper(log, util.Insecure())

	if (resp.Auth || resp.AuthEn) && (user == "" || password == "") {
		return conn, fmt.Errorf("%s (%s) missing user/password", resp.Model, resp.Mac)
	}

	switch conn.Gen {
	case 0, 1:
		// Shelly GEN 1 API
		// https://shelly-api-docs.shelly.cloud/gen1/#shelly-family-overview
		conn.URI = request.DefaultScheme(uri, "http")
		if user != "" {
			log.Redact(request.BasicAuthHeader(user, password))
			conn.Client.Transport = request.BasicAuth(user, password, conn.Client.Transport)
		}

	case 2:
		// Shelly GEN 2 API
		// https://shelly-api-docs.shelly.cloud/gen2/
		conn.URI = fmt.Sprintf("%s/rpc", request.DefaultScheme(uri, "http"))
		if user != "" {
			conn.Client.Transport = digest.NewTransport(user, password, conn.Client.Transport)
		}

	default:
		return conn, fmt.Errorf("%s (%s) unknown api generation (%d)", resp.Type, resp.Model, conn.Gen)
	}

	return conn, nil
}

// ExecGen2Cmd executes a shelly api gen1/gen2 command and provides the response
func (d *Connection) ExecGen2Cmd(method string, enable bool, res interface{}) error {
	// Shelly gen 2 rfc7616 authentication
	// https://shelly-api-docs.shelly.cloud/gen2/Overview/CommonDeviceTraits#authentication
	// https://datatracker.ietf.org/doc/html/rfc7616

	data := &util.Gen2RpcPost{
		ID:     d.Channel,
		On:     enable,
		Src:    "evcc",
		Method: method,
	}

	req, err := request.New(http.MethodPost, fmt.Sprintf("%s/%s", d.URI, method), request.MarshalJSON(data), request.JSONEncoding)
	if err != nil {
		return err
	}

	return d.DoJSON(req, &res)
}
