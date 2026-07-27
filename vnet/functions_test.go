package vnet

import (
	"errors"
	"net"
	"testing"
)

func TestNetInterfaces(t *testing.T) {
	origNetInterfaces = func() ([]net.Interface, error) {
		return []net.Interface{
			{},
		}, nil
	}
	defer func() {
		origNetInterfaces = net.Interfaces
	}()
	netInterfaces()
}

func TestNetInterfacesError(t *testing.T) {
	origNetInterfaces = func() ([]net.Interface, error) {
		return nil, errors.New("test")
	}
	defer func() {
		origNetInterfaces = net.Interfaces
	}()
	netInterfaces()
}
