package vnet

import (
	"net"

	"github.com/nitroshare/gomdns/list"
	"github.com/nitroshare/gomdns/syncpoint"
)

var (
	syncInit   = syncpoint.New()
	syncQueued = syncpoint.New()
	syncRead   = syncpoint.New()
)

type mockUDPConn struct {
	iface        *MockInterface
	chanRead     chan any
	chanReadRet  chan *Packet
	chanWrite    chan *Packet
	chanWriteRet chan any
	chanClose    chan any
	chanClosed   chan any
}

func (m *mockUDPConn) run() {
	defer close(m.chanClosed)
	var (
		chanQueued     = m.iface.queued.Subscribe()
		readQueue      = &list.List[*Packet]{}
		writeQueue     = &list.List[*Packet]{}
		waitingForRead bool
	)
	syncInit.Trigger()
	for {
		var (
			chanDequeue chan *Packet
			e           = writeQueue.PopFront()
			writePacket *Packet
		)
		if e != nil {
			chanDequeue = m.iface.chanDequeue
			writePacket = e.Value
		}
		select {
		case p := <-chanQueued:
			if waitingForRead {
				m.chanReadRet <- p
				waitingForRead = false
			} else {
				syncQueued.Trigger()
				readQueue.Add(p)
			}
		case <-m.chanRead:
			e := readQueue.PopFront()
			if e != nil {
				m.chanReadRet <- e.Value
			} else {
				syncRead.Trigger()
				waitingForRead = true
			}
		case p := <-m.chanWrite:
			writeQueue.Add(p)
			m.chanWriteRet <- nil
		case chanDequeue <- writePacket:
			writeQueue.PopFront()
		case <-m.chanClose:
			return
		}
	}
}

func newMockUDPConn(iface *MockInterface) *mockUDPConn {
	m := &mockUDPConn{
		iface:        iface,
		chanRead:     make(chan any),
		chanReadRet:  make(chan *Packet),
		chanWrite:    make(chan *Packet),
		chanWriteRet: make(chan any),
		chanClose:    make(chan any),
		chanClosed:   make(chan any),
	}
	go m.run()
	return m
}

func (m *mockUDPConn) ReadFrom(b []byte) (int, net.Addr, error) {
	select {
	case m.chanRead <- nil:
		select {
		case p := <-m.chanReadRet:
			return copy(b, p.Data), p.Addr, nil
		case <-m.chanClose:
			return 0, nil, net.ErrClosed
		}
	case <-m.chanClose:
		return 0, nil, net.ErrClosed
	}
}

func (m *mockUDPConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	p := &Packet{
		Addr: addr,
		Data: b,
	}
	select {
	case m.chanWrite <- p:
		<-m.chanWriteRet
		return len(b), nil
	case <-m.chanClose:
		return 0, net.ErrClosed
	}
}

func (m *mockUDPConn) Close() error {
	close(m.chanClose)
	<-m.chanClosed
	return nil
}
