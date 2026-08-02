package vnet

import (
	"net"
	"testing"

	"github.com/nitroshare/gomdns/compare"
)

var (
	testUDPAddr net.Addr = nil
	testUDPData          = []byte("test")
	testPacket           = &Packet{
		Addr: testUDPAddr,
		Data: testUDPData,
	}
)

func TestUDPConn(t *testing.T) {
	i := NewMockInterface()
	defer i.Close()

	t.Run("Read with data ready", func(t *testing.T) {
		syncQueued.Activate()
		defer syncQueued.Deactivate()
		m, _ := i.Listen("", nil)
		defer m.Close()
		i.EnqueueRead(testPacket)
		syncQueued.Wait()
		b := make([]byte, 32)
		n, addr, err := m.ReadFrom(b)
		compare.Compare(t, string(b[:n]), string(testUDPData), true)
		compare.Compare(t, addr, testUDPAddr, true)
		compare.Compare(t, err, nil, true)
	})

	t.Run("Blocking read", func(t *testing.T) {
		syncRead.Activate()
		defer syncRead.Deactivate()
		m, _ := i.Listen("", nil)
		defer m.Close()
		var (
			chanClose = make(chan any)
			b         = make([]byte, 32)
			n         int
			addr      net.Addr
			err       error
		)
		go func() {
			defer close(chanClose)
			n, addr, err = m.ReadFrom(b)
		}()
		syncRead.Wait()
		i.EnqueueRead(testPacket)
		<-chanClose
		compare.Compare(t, string(b[:n]), string(testUDPData), true)
		compare.Compare(t, addr, testUDPAddr, true)
		compare.Compare(t, err, nil, true)
	})

	t.Run("Reading when closing and closed", func(t *testing.T) {
		syncRead.Activate()
		defer syncRead.Deactivate()
		m, _ := i.Listen("", nil)
		var (
			chanClose = make(chan any)
			b         = make([]byte, 32)
			err       error
		)
		go func() {
			defer close(chanClose)
			_, _, err = m.ReadFrom(b)
		}()
		syncRead.Wait()
		m.Close()
		<-chanClose
		compare.Compare(t, err, net.ErrClosed, true)
		_, _, err = m.ReadFrom(nil)
		compare.Compare(t, err, net.ErrClosed, true)
	})

	t.Run("Dequeue", func(t *testing.T) {
		m, _ := i.Listen("", nil)
		defer m.Close()
		m.WriteTo(testUDPData, testUDPAddr)
		p, err := i.DequeueWrite()
		compare.Compare(t, string(p.Data), string(testUDPData), true)
		compare.Compare(t, p.Addr, testUDPAddr, true)
		compare.Compare(t, err, nil, true)
	})

	t.Run("Writing when closed", func(t *testing.T) {
		m, _ := i.Listen("", nil)
		m.Close()
		_, err := m.WriteTo(nil, nil)
		compare.Compare(t, err, net.ErrClosed, true)
	})
}
