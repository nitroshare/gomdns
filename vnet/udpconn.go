package vnet

import (
	"net"
)

// UDPConn provides the bare minimum methods required to interact with a UDP
// socket for sending and receiving datagrams.
type UDPConn interface {
	ReadFrom([]byte) (int, net.Addr, error)
	WriteTo([]byte, net.Addr) (int, error)
	Close() error
}
