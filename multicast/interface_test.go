package multicast

import (
	"errors"
	"net"
	"testing"

	"github.com/nitroshare/gomdns/compare"
)

func TestNetInterfaceMethods(t *testing.T) {
	i := netInterface{
		i: &net.Interface{},
	}

	// Invoke these, but we can't really inspect the return values since the
	// net.Interface we constructed is invalid

	i.Addrs()
	i.Flags()
	i.Interface()
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

func TestMockInterfacesWithError(t *testing.T) {
	MockWithError()
	defer Unmock()
	i, err := Interfaces()
	compare.Compare(t, i == nil, true, true)
	compare.Compare(t, err != nil, true, true)
}

func TestAddClearMockInterface(t *testing.T) {
	i, _ := mockInterfaces()
	compare.Compare(t, len(i), 0, true)
	AddMockInterface(&MockInterface{})
	i, _ = mockInterfaces()
	compare.Compare(t, len(i), 1, true)
	ClearMockInterfaces()
	i, _ = mockInterfaces()
	compare.Compare(t, len(i), 0, true)
}
