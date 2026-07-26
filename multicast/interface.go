package multicast

import (
	"errors"
	"net"
)

// Interface wraps the information in net.Interface. This allows for much
// easier mocking.
type Interface interface {

	// Addrs returns the addresses on the interface.
	Addrs() ([]net.Addr, error)

	// Flags returns the flags on the interface.
	Flags() net.Flags

	// Interface returns the underlying net.Interface.
	Interface() *net.Interface
}

type netInterface struct {
	i *net.Interface
}

func (n netInterface) Addrs() ([]net.Addr, error) { return n.i.Addrs() }
func (n netInterface) Flags() net.Flags           { return n.i.Flags }
func (n netInterface) Interface() *net.Interface  { return n.i }

// Packet represents a packet sent or received.
type Packet struct {
	Addr net.Addr
	Data []byte
}

var (
	// Interfaces returns a list of interfaces on the host when not mocked.
	Interfaces func() ([]Interface, error)

	origNetInterfaces = net.Interfaces
)

func netInterfaces() ([]Interface, error) {
	l, err := origNetInterfaces()
	if err != nil {
		return nil, err
	}
	interfaces := []Interface{}
	for _, i := range l {
		interfaces = append(interfaces, &netInterface{&i})
	}
	return interfaces, nil
}

var (
	mInterfaces []Interface
)

func mockInterfaces() ([]Interface, error) {
	return mInterfaces, nil
}

func mockInterfacesWithError() ([]Interface, error) {
	return nil, errors.New("mocked error")
}

// AddMockInterface adds a MockInterface to the list returned by Interfaces
// when mocking is enabled.
func AddMockInterface(i *MockInterface) {
	mInterfaces = append(mInterfaces, i)
}

// ClearMockInterfaces removes all mocked interfaces.
func ClearMockInterfaces() {
	mInterfaces = []Interface{}
}
