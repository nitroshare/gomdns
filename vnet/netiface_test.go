package vnet

import (
	"net"
	"testing"
)

func TestNetInterface(t *testing.T) {
	origNetListenMulticastUDP = func(string, *net.Interface, *net.UDPAddr) (*net.UDPConn, error) {
		return nil, nil
	}
	defer func() {
		origNetListenMulticastUDP = net.ListenMulticastUDP
	}()
	n := netInterface{
		i: &net.Interface{},
	}
	n.Name()
	n.Addrs()
	n.Flags()
	n.Listen("udp4", nil)
}
