# Shelly Client
[![Build and Test](https://github.com/darkmatus/shelly/actions/workflows/build.yml/badge.svg)](https://github.com/darkmatus/shelly/actions/workflows/build.yml)

This is a shelly api client for golang.
It provides functions for shelly's with and without energymeter.

## Implemented Shelly's

* Shelly 1Plus

## Usage
Call `go get github.com/darkmatus/shelly` in your project.

Create your required instance (current only Shelly 1plus is supported):

```go
    client, err := shelly.NewShelly(deviceType, authKey, baseURL, deviceID)
    if err != nil {
        return
    }
```
Use the `shelly.Device*` constants to create your device.
To use the Cloud Control API, you need the API key and the server uri.
Both can be found in the shelly cloud app or in the [shelly control center](https://control.shelly.cloud/).

## How to add a client.
If you want to add a new client, implement the ShellyInterface for devices which are switches. The interface requires
`on`, `off`, `toggle` and `getDeviceStatus`.
Don't forget to write tests. Add create a merge request.
