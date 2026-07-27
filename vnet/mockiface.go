package vnet

import (
	"net"

	"github.com/nitroshare/gomdns/broadcaster"
)

// MockInterface provides a fully mocked Interface that can be used for
// creating UDPConn-compatible connections for sending and receiving
// datagrams.
type MockInterface struct {

	// MockAddrs provides a list of addresses to be returned by Addrs.
	MockAddrs []net.Addr

	// MockAddrsError provides an error to be returned by Addrs.
	MockAddrsError error

	// MockFlags provides flags to be returned by Flags()
	MockFlags net.Flags

	queued      *broadcaster.Broadcaster[*Packet]
	chanDequeue chan *Packet
}

// NewMockInterface creates a new MockInterface.
func NewMockInterface() *MockInterface {
	return &MockInterface{
		queued:      broadcaster.New[*Packet](),
		chanDequeue: make(chan *Packet),
	}
}

func (m *MockInterface) Name() string               { return "MockInterface" }
func (m *MockInterface) Addrs() ([]net.Addr, error) { return m.MockAddrs, m.MockAddrsError }
func (m *MockInterface) Flags() net.Flags           { return m.MockFlags }

func (m *MockInterface) Listen(string, *net.UDPAddr) (UDPConn, error) {
	return newMockUDPConn(m), nil
}

// EnqueueRead queues a packet for ReadFrom in all UDPConn instances.
func (m *MockInterface) EnqueueRead(p *Packet) {
	m.queued.Send(p)
}

// DequeueWrite dequeues the next available packet written to a UDPConn if
// available, waiting until one is available.
func (m *MockInterface) DequeueWrite() (*Packet, error) {
	return <-m.chanDequeue, nil
}

// Close shuts down the MockInterface.
func (m *MockInterface) Close() {
	m.queued.Close()
}
