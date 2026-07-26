package multicast

import (
	"errors"
	"net"
	"testing"

	"github.com/nitroshare/gomdns/compare"
)

func TestNetInterface(t *testing.T) {
	compare.Compare(t, netInterface{}.Interface(), nil, true)
}

func TestInterfaces(t *testing.T) {
	_, err := Interfaces()
	compare.Compare(t, err, nil, true)
}

func TestInterfacesError(t *testing.T) {
	errTest := errors.New("test")
	origNetInterfaces = func() ([]net.Interface, error) { return nil, errTest }
	defer func() { origNetInterfaces = net.Interfaces }()
	_, err := Interfaces()
	compare.Compare(t, err, errTest, true)
}

func TestAddMockedInterface(t *testing.T) {
	i, _ := mockInterfaces()
	compare.Compare(t, len(i), 0, true)
	AddMockInterface(&MockInterface{})
	i, _ = mockInterfaces()
	compare.Compare(t, len(i), 1, true)
}
