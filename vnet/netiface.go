package vnet

import (
	"net"
)

var (
	origNetListenMulticastUDP = net.ListenMulticastUDP
)

type netInterface struct {
	i *net.Interface
}

func (n *netInterface) Name() string               { return n.i.Name }
func (n *netInterface) Addrs() ([]net.Addr, error) { return n.i.Addrs() }
func (n *netInterface) Flags() net.Flags           { return n.i.Flags }

func (n *netInterface) Listen(network string, gaddr *net.UDPAddr) (UDPConn, error) {
	return origNetListenMulticastUDP(network, n.i, gaddr)
}
