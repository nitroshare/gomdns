package watcher

import (
	"time"

	"github.com/nitroshare/gomdns/vnet"
)

// Config provides configuration for Watcher.
type Config struct {

	// Interval specifies how often the adapters on the host should be
	// enumerated.
	Interval time.Duration

	// ChanAdded receives an Interface when a new interface is added. This
	// cannot be nil and must be a valid channel.
	ChanAdded chan<- vnet.Interface

	// ChanRemoved receives an Interface when an interface is removed. This
	// cannot be nil and must be a valid channel.
	ChanRemoved chan<- vnet.Interface
}
