package multicast

import (
	"net"

	"github.com/nitroshare/gomdns/list"
)

// MockInterface implements Interface and provides methods for sending and
// receiving packets on the interface. Only one Listener should be created
// for each MockInterface.
type MockInterface struct {

	// MockAddrs contains a list of mocked addresses.
	MockAddrs []net.Addr

	// MockFlags contains mocked flags.
	MockFlags net.Flags

	chanQueueForRead    chan *Packet
	chanQueueForReadRet chan any
	chanReadRequest     chan any
	chanReadReply       chan *Packet
	chanWrite           chan *Packet
	chanWriteRet        chan any
	chanDequeueRequest  chan any
	chanDequeueReply    chan *Packet
	chanTest            chan any
	chanClose           chan any
	chanClosed          chan any
}

func (m *MockInterface) run() {
	defer close(m.chanClosed)
	defer close(m.chanReadReply)
	defer close(m.chanDequeueReply)
	var (
		readQueue         = &list.List[*Packet]{}
		writeQueue        = &list.List[*Packet]{}
		waitingForRead    bool
		waitingForDequeue bool
	)
	for {
		select {
		case p := <-m.chanQueueForRead:
			if waitingForRead {
				m.chanReadReply <- p
				waitingForRead = false
			} else {
				readQueue.Add(p)
			}
			m.chanQueueForReadRet <- nil
		case <-m.chanReadRequest:
			e := readQueue.PopFront()
			if e != nil {
				m.chanReadReply <- e.Value
			} else {
				if m.chanTest != nil {
					m.chanTest <- nil
				}
				waitingForRead = true
			}
		case p := <-m.chanWrite:
			if waitingForDequeue {
				m.chanDequeueReply <- p
				waitingForDequeue = false
			} else {
				writeQueue.Add(p)
			}
			m.chanWriteRet <- nil
		case <-m.chanDequeueRequest:
			e := writeQueue.PopFront()
			if e != nil {
				m.chanDequeueReply <- e.Value
			} else {
				if m.chanTest != nil {
					m.chanTest <- nil
				}
				waitingForDequeue = true
			}
		case <-m.chanClose:
			return
		}
	}
}

// NewMockInterface initializes a new instance of MockInterface.
func NewMockInterface() *MockInterface {
	m := &MockInterface{
		chanQueueForRead:    make(chan *Packet),
		chanQueueForReadRet: make(chan any),
		chanReadRequest:     make(chan any),
		chanReadReply:       make(chan *Packet),
		chanWrite:           make(chan *Packet),
		chanWriteRet:        make(chan any),
		chanDequeueRequest:  make(chan any),
		chanDequeueReply:    make(chan *Packet),
		chanClose:           make(chan any),
		chanClosed:          make(chan any),
	}
	go m.run()
	return m
}

func (m *MockInterface) Addrs() ([]net.Addr, error) {
	return m.MockAddrs, nil
}

func (m *MockInterface) Flags() net.Flags {
	return m.MockFlags
}

func (m *MockInterface) Interface() *net.Interface {
	return &net.Interface{
		MTU: 1500,
	}
}

// QueueForRead queues the provided packet for reading in a subsequent call to
// Listener.Read().
func (m *MockInterface) QueueForRead(p *Packet) {
	m.chanQueueForRead <- p
	<-m.chanQueueForReadRet
}

// DequeueWrite returns packets written to the Listener in order. This call
// will block until a packet is written or the connection is closed.
func (m *MockInterface) DequeueWrite() (*Packet, error) {
	m.chanDequeueRequest <- nil
	p, ok := <-m.chanDequeueReply
	if !ok {
		return nil, net.ErrClosed
	}
	return p, nil
}

type mockUDPConn struct {
	i *MockInterface
}

func (m *mockUDPConn) ReadFrom(b []byte) (int, net.Addr, error) {
	select {
	case _, ok := <-m.i.chanClose:
		if !ok {
			return 0, nil, net.ErrClosed
		}
	default:
	}
	m.i.chanReadRequest <- nil
	p, ok := <-m.i.chanReadReply
	if !ok {
		return 0, nil, net.ErrClosed
	}
	return copy(b, p.Data), p.Addr, nil
}

func (m *mockUDPConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	m.i.chanWrite <- &Packet{
		Addr: addr,
		Data: b,
	}
	<-m.i.chanWriteRet
	return len(b), nil
}

func (m *mockUDPConn) Close() error {
	close(m.i.chanClose)
	<-m.i.chanClosed
	return nil
}
