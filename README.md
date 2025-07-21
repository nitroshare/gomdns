## gomdns

[![Go Reference](https://pkg.go.dev/badge/github.com/nitroshare/gomdns.svg)](https://pkg.go.dev/github.com/nitroshare/gomdns)
[![MIT License](https://img.shields.io/badge/license-MIT-9370d8.svg?style=flat)](https://opensource.org/licenses/MIT)

This package aims to provide an [RFC 6762](https://datatracker.ietf.org/doc/html/rfc6762) compliant mDNS package for Go applications, with a heavy focus on simplicity. Although there are existing mDNS packages for Go, each of them lacked something we wanted, leading to the start of this package.

### Features

- Browser for continuously monitoring other devices providing a service
- Provider for exposing a local service on the network
- Ability to easily change parameters without recreating everything
- Comprehensive test suite to ensure compliance

This package is heavily based on [QMdnsEngine](https://github.com/nitroshare/qmdnsengine).

### Browser Example

Want to find devices on the network that provide `_http._tcp`?

```golang
import "github.com/nitroshare/gomdns/browser"

// Channels receive *Device when a device is added or removed
var (
    chanAdded   = make(chan *Device)
    chanRemoved = make(chan *Device)
)

// Create the browser
b, _ := browser.New(&browser.Config{
    Service:     "_http._tcp",
    ChanAdded:   chanAdded,
    ChanRemoved: chanRemoved,
})

// Read from chanAdded or chanRemoved in a separate goroutine
d := <-chanAdded
d := <-chanRemoved

// Close the browser when you are done
b.Close()
```
