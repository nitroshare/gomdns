package vnet

import (
	"net"
)

// Interface provides methods for working with interfaces, either directly
// with net.Interface or through mocked interfaces.
type Interface interface {

	// Addrs returns the addresses on the interface.
	Addrs() ([]net.Addr, error)

	// Flags returns the flags on the interface.
	Flags() net.Flags

	// Interface returns the underlying net.Interface.
	Interface() *net.Interface
}
