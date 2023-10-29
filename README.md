# Shelly Client
[![Build and Test](https://github.com/darkmatus/shelly/actions/workflows/build.yml/badge.svg)](https://github.com/darkmatus/shelly/actions/workflows/build.yml)

# This does "must" not work, actually. The work is not finished and untested!
This is a shelly api client for golang.
It provides functions for shelly's with and without energymeter.

## Usage
Call `go get github.com/darkmatus/shelly` in your project.

Create your required instance (current only Shelly 1plus is supported):

```go
	client, err := shelly.NewShelly(deviceType, authKey, baseURL, deviceID)
    if err != nil {
        return
    }
```
Use the shelly.Device* constants to create your device.

## How to add a client.
If you want to add a new client, implement the ShellyInterface for devices which are switches. The interface requires
`on`, `off`, `toggle` and `getDeviceStatus`.
Don't forget to write tests. Add create a merge request.
