package multicast

import (
	"net"
	"testing"

	"github.com/nitroshare/gomdns/compare"
)

var (
	testAddr net.Addr = nil
	testData          = []byte("test")
)

func TestMockInterfaceMethods(t *testing.T) {
	i := NewMockInterface()
	_, err := i.Addrs()
	if err != nil {
		t.Fatal(err)
	}
	i.Flags()
	i.Interface()
}

func TestMockInterfaceQueueForRead(t *testing.T) {
	var (
		m = NewMockInterface()
		c = &mockUDPConn{m}
	)
	m.chanTest = make(chan any)
	t.Run("read with data ready", func(t *testing.T) {
		var (
			b = make([]byte, 32)
		)
		m.QueueForRead(&Packet{testAddr, testData})
		n, addr, err := c.ReadFrom(b)
		compare.Compare(t, string(b[:n]), string(testData), true)
		compare.Compare(t, n, len(testData), true)
		compare.Compare(t, addr, testAddr, true)
		compare.Compare(t, err, nil, true)
	})
	t.Run("blocking read", func(t *testing.T) {
		var (
			chanClose = make(chan any)
			b         = make([]byte, 32)
			n         int
			addr      net.Addr
			err       error
		)
		go func() {
			defer close(chanClose)
			n, addr, err = c.ReadFrom(b)
		}()
		<-m.chanTest
		m.QueueForRead(&Packet{testAddr, testData})
		<-chanClose
		compare.Compare(t, string(b[:n]), string(testData), true)
		compare.Compare(t, n, len(testData), true)
		compare.Compare(t, addr, testAddr, true)
		compare.Compare(t, err, nil, true)
	})
	t.Run("close while reading", func(t *testing.T) {
		var (
			chanClose = make(chan any)
			b         = make([]byte, 32)
			err       error
		)
		go func() {
			defer close(chanClose)
			_, _, err = c.ReadFrom(b)
		}()
		<-m.chanTest
		c.Close()
		<-chanClose
		compare.Compare(t, err, net.ErrClosed, true)
	})
}

func TestMockInterfaceDequeueWrite(t *testing.T) {
	var (
		m = NewMockInterface()
		c = &mockUDPConn{m}
	)
	m.chanTest = make(chan any)
	t.Run("dequeue with data ready", func(t *testing.T) {
		c.WriteTo(testData, testAddr)
		p, err := m.DequeueWrite()
		compare.Compare(t, p.Addr, testAddr, true)
		compare.Compare(t, string(p.Data), string(testData), true)
		compare.Compare(t, err, nil, true)
	})
	t.Run("blocking dequeue", func(t *testing.T) {
		var (
			chanClose = make(chan any)
			p         *Packet
			err       error
		)
		go func() {
			defer close(chanClose)
			p, err = m.DequeueWrite()
		}()
		<-m.chanTest
		n, err := c.WriteTo(testData, testAddr)
		compare.Compare(t, n, len(testData), true)
		compare.Compare(t, err, nil, true)
		<-chanClose
		compare.Compare(t, p.Addr, testAddr, true)
		compare.Compare(t, string(p.Data), string(testData), true)
		compare.Compare(t, err, nil, true)
	})
	t.Run("close while dequeing", func(t *testing.T) {
		var (
			chanClose = make(chan any)
			err       error
		)
		go func() {
			defer close(chanClose)
			_, err = m.DequeueWrite()
		}()
		<-m.chanTest
		c.Close()
		<-chanClose
		compare.Compare(t, err, net.ErrClosed, true)
	})
}
