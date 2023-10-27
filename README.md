# Shelly Client
[![Build and Test](https://github.com/darkmatus/shelly/actions/workflows/build.yml/badge.svg)](https://github.com/darkmatus/shelly/actions/workflows/build.yml)

# This does "must" not work, actually. The work is not finished and untested!
This is a shelly api client for golang.
It provides functions for shelly's with and without energymeter.

## Usage
Call `go get github.com/darkmatus/shelly` in your project.

Create your required instance:
```go
	var shelly, _ = NewShelly("123", "user", "pw", 0)
err := shelly.Enable(true)
if err != nil {
return
}
```
