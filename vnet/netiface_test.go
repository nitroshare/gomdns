package vnet

import (
	"net"
	"testing"
)

func TestNetInterface(t *testing.T) {
	n := netInterface{
		i: &net.Interface{},
	}
	n.Addrs()
	n.Flags()
	n.Interface()
}
