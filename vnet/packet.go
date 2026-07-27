package vnet

import (
	"net"
)

// Packet represents a datagram to be sent or received.
type Packet struct {

	// Addr indicates the packet's source or destination.
	Addr net.Addr

	// Data contains the contents of the packet.
	Data []byte
}
